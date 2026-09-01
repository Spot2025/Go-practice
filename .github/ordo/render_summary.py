#!/usr/bin/env python3
"""Render Ordo's stable JSON report as a GitHub-flavored Markdown summary."""

from __future__ import annotations

import argparse
import hashlib
import json
import os
import re
import sys
from pathlib import Path
from typing import Any, Mapping, Sequence
from urllib.parse import quote


_REPOSITORY_RE = re.compile(r"^[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+$")
_COMMIT_SHA_RE = re.compile(r"^[0-9A-Fa-f]{7,64}$")


class ReportError(ValueError):
    """Raised when input is not a valid Ordo JSON report."""


def load_report(path: Path | str) -> Mapping[str, Any]:
    """Load an Ordo report from *path* and return its decoded JSON object."""

    report_path = Path(path)
    try:
        with report_path.open(encoding="utf-8") as report_file:
            report = json.load(report_file)
    except OSError as error:
        raise ReportError(f"could not read {report_path}: {error}") from error
    except json.JSONDecodeError as error:
        raise ReportError(
            f"invalid JSON in {report_path} at line {error.lineno}, "
            f"column {error.colno}: {error.msg}"
        ) from error

    _validate_report(report)
    return report


def display_name(function_id: str) -> str:
    """Return a compact function name without trying to parse all Go syntax."""

    # SSA renders methods as "(*full/import/path.Type).Method" or
    # "(full/import/path.Type).Method". Keep pointer/value receiver notation,
    # but remove the import path from inside the receiver.
    method = re.fullmatch(r"\((\*)?(.+)\)\.([^/]+)", function_id)
    if method:
        pointer, qualified_type, method_name = method.groups()
        receiver = qualified_type.rsplit(".", 1)[-1]
        if pointer:
            return f"(*{receiver}).{method_name}()"
        return f"{receiver}.{method_name}()"

    # Also tolerate the more human-oriented form sometimes used in examples:
    # "full/import/path/pkg.(*Type).Method".
    receiver_suffix = re.search(r"\.\((\*?[^()]+)\)\.([^/]+)$", function_id)
    if receiver_suffix:
        receiver, method_name = receiver_suffix.groups()
        return f"({receiver}).{method_name}()"

    tail = function_id.rsplit("/", 1)[-1]
    parts = tail.split(".")
    if len(parts) >= 3:
        # Accommodate conceptual IDs such as "pkg.Type.Method" while keeping
        # this intentionally much smaller than a Go parser.
        short_name = ".".join(parts[-2:])
    elif len(parts) == 2:
        short_name = parts[-1]
    else:
        short_name = tail
    return f"{short_name}()"


def render_summary(
    report: Mapping[str, Any],
    repository: str | None = None,
    sha: str | None = None,
    source_prefix: str | None = None,
    pull_request: int | str | None = None,
) -> str:
    """Render a validated Ordo report as Markdown, including a final newline."""

    _validate_report(report)
    steps = report["steps"]
    lines = ["# Ordo — Review Path", ""]

    if not steps:
        lines.append("No changed Go functions were detected in this pull request.")
        return "\n".join(lines) + "\n"

    function_count = sum(len(step["functions"]) for step in steps)
    lines.extend(
        [
            f"**Changed functions:** {function_count}  ",
            f"**Review steps:** {len(steps)}",
            "",
        ]
    )

    full_ids: list[tuple[int, str]] = []
    for step in steps:
        order = step["order"]
        functions = step["functions"]
        cycle = step["cycle"]
        artifact = step.get("artifact", False)

        if cycle:
            lines.append(f"## {order} · Cycle · {len(functions)} functions")
            lines.append("")
            lines.append(
                "> 🔁 These functions form a dependency cycle; review them together."
            )
            if artifact:
                lines.append("")
                lines.append(
                    "> ⚠️ Ordo classified this cycle as a probable call-graph artifact."
                )
            lines.append("")
            for function in functions:
                name = _inline_code(display_name(function["id"]))
                location = _render_location(
                    function, repository, sha, source_prefix, pull_request
                )
                lines.append(f"- {name} — {location}")
                full_ids.append((order, function["id"]))
        else:
            function = functions[0]
            location = _render_location(
                function, repository, sha, source_prefix, pull_request
            )
            lines.append(f"## {order} · {location}")
            lines.append("")
            lines.append(_inline_code(display_name(function["id"])))
            full_ids.append((order, function["id"]))
        lines.append("")

    lines.extend(["<details>", "<summary>Full function IDs</summary>", ""])
    for order, function_id in full_ids:
        lines.append(f"- Step {order} — {_inline_code(function_id)}")
    lines.extend(["", "</details>"])
    return "\n".join(lines) + "\n"


def _render_location(
    function: Mapping[str, Any],
    repository: str | None,
    sha: str | None,
    source_prefix: str | None,
    pull_request: int | str | None,
) -> str:
    file_name = function["file"]
    line = function["line"]

    if not file_name:
        return "Source unavailable"

    repository_path, prefix_is_safe = _repository_path(file_name, source_prefix)
    location = f"{repository_path}:{line}" if line > 0 else repository_path
    location_markdown = _inline_code(location)
    source_url = (
        _source_url(repository, sha, repository_path, line)
        if prefix_is_safe
        else None
    )
    diff_url = (
        _pull_request_diff_url(
            repository, pull_request, repository_path, line
        )
        if prefix_is_safe
        else None
    )

    if diff_url is not None:
        diff_location = f"[{location_markdown}]({diff_url})"
        if source_url is not None:
            return f"{diff_location} ([source]({source_url}))"
        return diff_location
    if source_url is not None:
        return f"[{location_markdown}]({source_url})"
    return location_markdown


def _repository_path(file_name: str, source_prefix: str | None) -> tuple[str, bool]:
    prefix = _normalize_source_prefix(source_prefix)
    if prefix is None:
        return file_name, False
    if not prefix:
        return file_name, True
    return f"{prefix}/{file_name}", True


def _normalize_source_prefix(source_prefix: str | None) -> str | None:
    if source_prefix is None:
        return ""
    if source_prefix.startswith(("/", "\\")):
        return None
    if any(character in source_prefix for character in ("\x00", "\n", "\r", "\\")):
        return None

    components = source_prefix.split("/")
    while components and components[0] == ".":
        components.pop(0)
    while components and components[-1] == "":
        components.pop()
    if any(component in {"", ".", ".."} for component in components):
        return None
    return "/".join(components)


def _source_url(
    repository: str | None,
    sha: str | None,
    file_name: str,
    line: int,
) -> str | None:
    if not repository or not _REPOSITORY_RE.fullmatch(repository):
        return None
    if not sha or not _COMMIT_SHA_RE.fullmatch(sha):
        return None
    if not _is_safe_repository_path(file_name):
        return None

    encoded_path = quote(file_name, safe="/-._~")
    url = f"https://github.com/{repository}/blob/{sha}/{encoded_path}"
    if line > 0:
        url += f"#L{line}"
    return url


def _pull_request_diff_url(
    repository: str | None,
    pull_request: int | str | None,
    file_name: str,
    line: int,
) -> str | None:
    if not repository or not _REPOSITORY_RE.fullmatch(repository):
        return None
    pull_request_number = _normalize_pull_request(pull_request)
    if pull_request_number is None or line <= 0:
        return None
    if not _is_safe_repository_path(file_name):
        return None

    path_hash = hashlib.sha256(file_name.encode("utf-8")).hexdigest()
    return (
        f"https://github.com/{repository}/pull/{pull_request_number}/files"
        f"#diff-{path_hash}R{line}"
    )


def _normalize_pull_request(pull_request: int | str | None) -> str | None:
    if isinstance(pull_request, bool) or pull_request is None:
        return None
    if isinstance(pull_request, int):
        return str(pull_request) if pull_request > 0 else None
    if isinstance(pull_request, str) and re.fullmatch(r"[1-9][0-9]*", pull_request):
        return pull_request
    return None


def _is_safe_repository_path(file_name: str) -> bool:
    if not file_name or file_name.startswith(("/", "\\")):
        return False
    if "\x00" in file_name or "\n" in file_name or "\r" in file_name:
        return False
    return not any(
        component in {"", ".", ".."} for component in file_name.split("/")
    )


def _inline_code(value: str) -> str:
    value = value.replace("\r", " ").replace("\n", " ")
    longest_run = max((len(run) for run in re.findall(r"`+", value)), default=0)
    fence = "`" * (longest_run + 1)
    padding = " " if value.startswith("`") or value.endswith("`") else ""
    return f"{fence}{padding}{value}{padding}{fence}"


def _validate_report(report: Any) -> None:
    if not isinstance(report, Mapping):
        raise ReportError("Ordo JSON root must be an object")
    if "steps" not in report or not isinstance(report["steps"], list):
        raise ReportError('Ordo JSON must contain a "steps" array')

    for step_index, step in enumerate(report["steps"]):
        context = f"steps[{step_index}]"
        if not isinstance(step, Mapping):
            raise ReportError(f"{context} must be an object")
        _require_integer(step, "order", context, minimum=1)
        _require_boolean(step, "cycle", context)
        if "artifact" in step:
            _require_boolean(step, "artifact", context)
        functions = step.get("functions")
        if not isinstance(functions, list) or not functions:
            raise ReportError(f'{context}.functions must be a non-empty array')
        if not step["cycle"] and len(functions) != 1:
            raise ReportError(
                f"{context} has multiple functions but is not marked as a cycle"
            )
        if step.get("artifact", False) and not step["cycle"]:
            raise ReportError(f"{context} is an artifact but is not marked as a cycle")

        for function_index, function in enumerate(functions):
            function_context = f"{context}.functions[{function_index}]"
            if not isinstance(function, Mapping):
                raise ReportError(f"{function_context} must be an object")
            for field in ("id", "file"):
                if field not in function or not isinstance(function[field], str):
                    raise ReportError(f"{function_context}.{field} must be a string")
            if not function["id"]:
                raise ReportError(f"{function_context}.id must not be empty")
            _require_integer(function, "line", function_context, minimum=0)


def _require_integer(
    value: Mapping[str, Any], field: str, context: str, minimum: int
) -> None:
    field_value = value.get(field)
    if (
        not isinstance(field_value, int)
        or isinstance(field_value, bool)
        or field_value < minimum
    ):
        raise ReportError(f"{context}.{field} must be an integer >= {minimum}")


def _require_boolean(value: Mapping[str, Any], field: str, context: str) -> None:
    if field not in value or not isinstance(value[field], bool):
        raise ReportError(f"{context}.{field} must be a boolean")


def _parse_args(argv: Sequence[str] | None) -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Render an Ordo JSON report as GitHub-flavored Markdown."
    )
    parser.add_argument("report", type=Path, help="path to Ordo's JSON report")
    parser.add_argument(
        "--repository",
        default=os.environ.get("GITHUB_REPOSITORY"),
        metavar="OWNER/REPO",
        help="GitHub repository for source links (default: GITHUB_REPOSITORY)",
    )
    parser.add_argument(
        "--sha",
        default=os.environ.get("GITHUB_SHA"),
        help="commit SHA for source links (default: GITHUB_SHA)",
    )
    parser.add_argument(
        "--source-prefix",
        help="repository-relative prefix for paths emitted from a module subdirectory",
    )
    parser.add_argument(
        "--pull-request",
        default=os.environ.get("PULL_REQUEST_NUMBER"),
        metavar="NUMBER",
        help="pull request number for diff links (default: PULL_REQUEST_NUMBER)",
    )
    return parser.parse_args(argv)


def main(argv: Sequence[str] | None = None) -> int:
    args = _parse_args(argv)
    try:
        report = load_report(args.report)
        markdown = render_summary(
            report,
            args.repository,
            args.sha,
            args.source_prefix,
            args.pull_request,
        )
    except ReportError as error:
        print(f"render_summary.py: error: {error}", file=sys.stderr)
        return 1

    sys.stdout.write(markdown)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
