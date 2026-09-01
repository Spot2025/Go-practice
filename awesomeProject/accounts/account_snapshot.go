package accounts

import "awesomeProject/accounts/models"

func (h *Handler) accountSnapshot(name string) (models.Account, bool) {
	h.guard.RLock()
	defer h.guard.RUnlock()

	account, ok := h.accounts[name]
	if !ok {
		return models.Account{}, false
	}

	return *account, true
}
