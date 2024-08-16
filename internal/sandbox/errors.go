package sandbox

import "github.com/sourcenetwork/acp_core/pkg/errors"

func newNewSandboxErr(err error) error {
	return errors.Wrap("NewSandbox failed", err)
}

func newListSandboxesErr(err error) error {
	return errors.Wrap("ListSandboxes failed", err)
}

func newSetStateErr(err error, handle uint64) error {
	return errors.Wrap("SetState failed", err, errors.Pair("handle", handle))
}

func newVerifyTheoremsErr(err error) error {
	return errors.Wrap("VerifyTheorems failed", err)
}

func newGetCatalogueErr(err error) error {
	return errors.Wrap("GetCatalogue failed", err)
}

func newRestoreScratchpadErr(err error) error {
	return errors.Wrap("RestoreScratchpad failed", err)
}
