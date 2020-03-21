package webdriver

import (
	"context"

	"github.com/mediabuyerbot/go-webdriver/pkg/w3c"
)

type selector struct {
	id       string
	strategy w3c.FindElementStrategy
}

type WebElement struct {
	elem w3c.WebElement
	ctx  context.Context
	q    selector
}

func (w WebElement) Attr(name string) (string, error) {
	return w.elem.GetAttribute(w.ctx, name)
}

func (w WebElement) PressNullKey() error {
	return w.elem.SendKeys(w.ctx, w3c.NullKey)
}

func (w WebElement) PressCancelKey() error {
	return w.elem.SendKeys(w.ctx, w3c.CancelKey)
}

func (w WebElement) PressHelpKey() error {
	return w.elem.SendKeys(w.ctx, w3c.HelpKey)
}

func (w WebElement) PressBackspaceKey() error {
	return w.elem.SendKeys(w.ctx, w3c.BackspaceKey)
}

func (w WebElement) PressTabKey() error {
	return w.elem.SendKeys(w.ctx, w3c.TabKey)
}

func (w WebElement) PressClearKey() error {
	return w.elem.SendKeys(w.ctx, w3c.ClearKey)
}

func (w WebElement) PressReturnKey() error {
	return w.elem.SendKeys(w.ctx, w3c.ReturnKey)
}

func (w WebElement) PressEnterKey() error {
	return w.elem.SendKeys(w.ctx, w3c.EnterKey)
}

func (w WebElement) PressShiftKey() error {
	return w.elem.SendKeys(w.ctx, w3c.ShiftKey)
}

func (w WebElement) PressControlKey() error {
	return w.elem.SendKeys(w.ctx, w3c.ControlKey)
}

func (w WebElement) PressAltKey() error {
	return w.elem.SendKeys(w.ctx, w3c.AltKey)
}

func (w WebElement) PressPauseKey() error {
	return w.elem.SendKeys(w.ctx, w3c.PauseKey)
}

func (w WebElement) PressEscapeKey() error {
	return w.elem.SendKeys(w.ctx, w3c.EscapeKey)
}

func (w WebElement) PressSpaceKey() error {
	return w.elem.SendKeys(w.ctx, w3c.SpaceKey)
}

func (w WebElement) PressPageUpKey() error {
	return w.elem.SendKeys(w.ctx, w3c.PageUpKey)
}

func (w WebElement) PressPageDownKey() error {
	return w.elem.SendKeys(w.ctx, w3c.PageDownKey)
}

func (w WebElement) PressEndKey() error {
	return w.elem.SendKeys(w.ctx, w3c.EndKey)
}

func (w WebElement) PressHomeKey() error {
	return w.elem.SendKeys(w.ctx, w3c.HomeKey)
}

func (w WebElement) PressLeftArrowKey() error {
	return w.elem.SendKeys(w.ctx, w3c.LeftArrowKey)
}

func (w WebElement) PressUpArrowKey() error {
	return w.elem.SendKeys(w.ctx, w3c.UpArrowKey)
}

func (w WebElement) PressRightArrowKey() error {
	return w.elem.SendKeys(w.ctx, w3c.RightArrowKey)
}

func (w WebElement) PressDownArrowKey() error {
	return w.elem.SendKeys(w.ctx, w3c.DownArrowKey)
}

func (w WebElement) PressInsertKey() error {
	return w.elem.SendKeys(w.ctx, w3c.InsertKey)
}

func (w WebElement) PressDeleteKey() error {
	return w.elem.SendKeys(w.ctx, w3c.DeleteKey)
}

func (w WebElement) PressSemicolonKey() error {
	return w.elem.SendKeys(w.ctx, w3c.SemicolonKey)
}

func (w WebElement) PressEqualsKey() error {
	return w.elem.SendKeys(w.ctx, w3c.EqualsKey)
}

func (w WebElement) PressNumpad0Key() error {
	return w.elem.SendKeys(w.ctx, w3c.Numpad0Key)
}

func (w WebElement) PressNumpad1Key() error {
	return w.elem.SendKeys(w.ctx, w3c.Numpad1Key)
}

func (w WebElement) PressNumpad2Key() error {
	return w.elem.SendKeys(w.ctx, w3c.Numpad2Key)
}

func (w WebElement) PressNumpad3Key() error {
	return w.elem.SendKeys(w.ctx, w3c.Numpad3Key)
}

func (w WebElement) PressNumpad4Key() error {
	return w.elem.SendKeys(w.ctx, w3c.Numpad4Key)
}

func (w WebElement) PressNumpad5Key() error {
	return w.elem.SendKeys(w.ctx, w3c.Numpad5Key)
}

func (w WebElement) PressNumpad6Key() error {
	return w.elem.SendKeys(w.ctx, w3c.Numpad6Key)
}

func (w WebElement) PressNumpad7Key() error {
	return w.elem.SendKeys(w.ctx, w3c.Numpad7Key)
}

func (w WebElement) PressNumpad8Key() error {
	return w.elem.SendKeys(w.ctx, w3c.Numpad8Key)
}

func (w WebElement) PressNumpad9Key() error {
	return w.elem.SendKeys(w.ctx, w3c.Numpad9Key)
}

func (w WebElement) PressMultiplyKey() error {
	return w.elem.SendKeys(w.ctx, w3c.MultiplyKey)
}

func (w WebElement) PressAddKey() error {
	return w.elem.SendKeys(w.ctx, w3c.AddKey)
}

func (w WebElement) PressSeparatorKey() error {
	return w.elem.SendKeys(w.ctx, w3c.SeparatorKey)
}

func (w WebElement) PressSubstractKey() error {
	return w.elem.SendKeys(w.ctx, w3c.SubstractKey)
}

func (w WebElement) PressDecimalKey() error {
	return w.elem.SendKeys(w.ctx, w3c.DecimalKey)
}

func (w WebElement) PressDivideKey() error {
	return w.elem.SendKeys(w.ctx, w3c.DivideKey)
}

func (w WebElement) PressF1Key() error {
	return w.elem.SendKeys(w.ctx, w3c.F1Key)
}

func (w WebElement) PressF2Key() error {
	return w.elem.SendKeys(w.ctx, w3c.F2Key)
}

func (w WebElement) PressF3Key() error {
	return w.elem.SendKeys(w.ctx, w3c.F3Key)
}

func (w WebElement) PressF4Key() error {
	return w.elem.SendKeys(w.ctx, w3c.F4Key)
}

func (w WebElement) PressF5Key() error {
	return w.elem.SendKeys(w.ctx, w3c.F5Key)
}

func (w WebElement) PressF6Key() error {
	return w.elem.SendKeys(w.ctx, w3c.F6Key)
}

func (w WebElement) PressF7Key() error {
	return w.elem.SendKeys(w.ctx, w3c.F7Key)
}

func (w WebElement) PressF8Key() error {
	return w.elem.SendKeys(w.ctx, w3c.F8Key)
}

func (w WebElement) PressF10Key() error {
	return w.elem.SendKeys(w.ctx, w3c.F9Key)
}

func (w WebElement) PressF11Key() error {
	return w.elem.SendKeys(w.ctx, w3c.F9Key)
}

func (w WebElement) PressF12Key() error {
	return w.elem.SendKeys(w.ctx, w3c.F9Key)
}

func (w WebElement) PressMetaKey() error {
	return w.elem.SendKeys(w.ctx, w3c.MetaKey)
}

func (w WebElement) SendKeys(keys ...w3c.Key) error {
	return w.elem.SendKeys(w.ctx, keys...)
}
