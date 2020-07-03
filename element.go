package webdriver

import (
	"context"

	"github.com/mediabuyerbot/go-webdriver/pkg/w3cproto"
)

type selector struct {
	id       string
	strategy w3cproto.FindElementStrategy
}

type WebElement struct {
	elem w3cproto.WebElement
	ctx  context.Context
	q    selector
}

func (w WebElement) Attr(name string) (string, error) {
	return w.elem.GetAttribute(w.ctx, name)
}

func (w WebElement) PressNullKey() error {
	return w.elem.SendKeys(w.ctx, w3cproto.NullKey)
}

func (w WebElement) PressCancelKey() error {
	return w.elem.SendKeys(w.ctx, w3cproto.CancelKey)
}

func (w WebElement) PressHelpKey() error {
	return w.elem.SendKeys(w.ctx, w3cproto.HelpKey)
}

func (w WebElement) PressBackspaceKey() error {
	return w.elem.SendKeys(w.ctx, w3cproto.BackspaceKey)
}

func (w WebElement) PressTabKey() error {
	return w.elem.SendKeys(w.ctx, w3cproto.TabKey)
}

func (w WebElement) PressClearKey() error {
	return w.elem.SendKeys(w.ctx, w3cproto.ClearKey)
}

func (w WebElement) PressReturnKey() error {
	return w.elem.SendKeys(w.ctx, w3cproto.ReturnKey)
}

func (w WebElement) PressEnterKey() error {
	return w.elem.SendKeys(w.ctx, w3cproto.EnterKey)
}

func (w WebElement) PressShiftKey() error {
	return w.elem.SendKeys(w.ctx, w3cproto.ShiftKey)
}

func (w WebElement) PressControlKey() error {
	return w.elem.SendKeys(w.ctx, w3cproto.ControlKey)
}

func (w WebElement) PressAltKey() error {
	return w.elem.SendKeys(w.ctx, w3cproto.AltKey)
}

func (w WebElement) PressPauseKey() error {
	return w.elem.SendKeys(w.ctx, w3cproto.PauseKey)
}

func (w WebElement) PressEscapeKey() error {
	return w.elem.SendKeys(w.ctx, w3cproto.EscapeKey)
}

func (w WebElement) PressSpaceKey() error {
	return w.elem.SendKeys(w.ctx, w3cproto.SpaceKey)
}

func (w WebElement) PressPageUpKey() error {
	return w.elem.SendKeys(w.ctx, w3cproto.PageUpKey)
}

func (w WebElement) PressPageDownKey() error {
	return w.elem.SendKeys(w.ctx, w3cproto.PageDownKey)
}

func (w WebElement) PressEndKey() error {
	return w.elem.SendKeys(w.ctx, w3cproto.EndKey)
}

func (w WebElement) PressHomeKey() error {
	return w.elem.SendKeys(w.ctx, w3cproto.HomeKey)
}

func (w WebElement) PressLeftArrowKey() error {
	return w.elem.SendKeys(w.ctx, w3cproto.LeftArrowKey)
}

func (w WebElement) PressUpArrowKey() error {
	return w.elem.SendKeys(w.ctx, w3cproto.UpArrowKey)
}

func (w WebElement) PressRightArrowKey() error {
	return w.elem.SendKeys(w.ctx, w3cproto.RightArrowKey)
}

func (w WebElement) PressDownArrowKey() error {
	return w.elem.SendKeys(w.ctx, w3cproto.DownArrowKey)
}

func (w WebElement) PressInsertKey() error {
	return w.elem.SendKeys(w.ctx, w3cproto.InsertKey)
}

func (w WebElement) PressDeleteKey() error {
	return w.elem.SendKeys(w.ctx, w3cproto.DeleteKey)
}

func (w WebElement) PressSemicolonKey() error {
	return w.elem.SendKeys(w.ctx, w3cproto.SemicolonKey)
}

func (w WebElement) PressEqualsKey() error {
	return w.elem.SendKeys(w.ctx, w3cproto.EqualsKey)
}

func (w WebElement) PressNumpad0Key() error {
	return w.elem.SendKeys(w.ctx, w3cproto.Numpad0Key)
}

func (w WebElement) PressNumpad1Key() error {
	return w.elem.SendKeys(w.ctx, w3cproto.Numpad1Key)
}

func (w WebElement) PressNumpad2Key() error {
	return w.elem.SendKeys(w.ctx, w3cproto.Numpad2Key)
}

func (w WebElement) PressNumpad3Key() error {
	return w.elem.SendKeys(w.ctx, w3cproto.Numpad3Key)
}

func (w WebElement) PressNumpad4Key() error {
	return w.elem.SendKeys(w.ctx, w3cproto.Numpad4Key)
}

func (w WebElement) PressNumpad5Key() error {
	return w.elem.SendKeys(w.ctx, w3cproto.Numpad5Key)
}

func (w WebElement) PressNumpad6Key() error {
	return w.elem.SendKeys(w.ctx, w3cproto.Numpad6Key)
}

func (w WebElement) PressNumpad7Key() error {
	return w.elem.SendKeys(w.ctx, w3cproto.Numpad7Key)
}

func (w WebElement) PressNumpad8Key() error {
	return w.elem.SendKeys(w.ctx, w3cproto.Numpad8Key)
}

func (w WebElement) PressNumpad9Key() error {
	return w.elem.SendKeys(w.ctx, w3cproto.Numpad9Key)
}

func (w WebElement) PressMultiplyKey() error {
	return w.elem.SendKeys(w.ctx, w3cproto.MultiplyKey)
}

func (w WebElement) PressAddKey() error {
	return w.elem.SendKeys(w.ctx, w3cproto.AddKey)
}

func (w WebElement) PressSeparatorKey() error {
	return w.elem.SendKeys(w.ctx, w3cproto.SeparatorKey)
}

func (w WebElement) PressSubstractKey() error {
	return w.elem.SendKeys(w.ctx, w3cproto.SubstractKey)
}

func (w WebElement) PressDecimalKey() error {
	return w.elem.SendKeys(w.ctx, w3cproto.DecimalKey)
}

func (w WebElement) PressDivideKey() error {
	return w.elem.SendKeys(w.ctx, w3cproto.DivideKey)
}

func (w WebElement) PressF1Key() error {
	return w.elem.SendKeys(w.ctx, w3cproto.F1Key)
}

func (w WebElement) PressF2Key() error {
	return w.elem.SendKeys(w.ctx, w3cproto.F2Key)
}

func (w WebElement) PressF3Key() error {
	return w.elem.SendKeys(w.ctx, w3cproto.F3Key)
}

func (w WebElement) PressF4Key() error {
	return w.elem.SendKeys(w.ctx, w3cproto.F4Key)
}

func (w WebElement) PressF5Key() error {
	return w.elem.SendKeys(w.ctx, w3cproto.F5Key)
}

func (w WebElement) PressF6Key() error {
	return w.elem.SendKeys(w.ctx, w3cproto.F6Key)
}

func (w WebElement) PressF7Key() error {
	return w.elem.SendKeys(w.ctx, w3cproto.F7Key)
}

func (w WebElement) PressF8Key() error {
	return w.elem.SendKeys(w.ctx, w3cproto.F8Key)
}

func (w WebElement) PressF10Key() error {
	return w.elem.SendKeys(w.ctx, w3cproto.F9Key)
}

func (w WebElement) PressF11Key() error {
	return w.elem.SendKeys(w.ctx, w3cproto.F9Key)
}

func (w WebElement) PressF12Key() error {
	return w.elem.SendKeys(w.ctx, w3cproto.F9Key)
}

func (w WebElement) PressMetaKey() error {
	return w.elem.SendKeys(w.ctx, w3cproto.MetaKey)
}

func (w WebElement) SendKeys(keys ...w3cproto.Key) error {
	return w.elem.SendKeys(w.ctx, keys...)
}
