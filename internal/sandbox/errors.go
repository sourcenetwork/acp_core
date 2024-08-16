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

func newVerifyTheoremsErr(err error, handle uint64) error {
	return errors.Wrap("VerifyTheorems failed", err, errors.Pair("handle", handle))
}

func newGetCatalogueErr(err error, handle uint64) error {
	return errors.Wrap("GetCatalogue failed", err, errors.Pair("handle", handle))
}

func newRestoreScratchpadErr(err error, handle uint64) error {
	return errors.Wrap("RestoreScratchpad failed", err, errors.Pair("handle", handle))
}
