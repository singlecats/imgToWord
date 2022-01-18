package doc

import (
	"baliance.com/gooxml/common"
	"baliance.com/gooxml/document"
	"baliance.com/gooxml/measurement"
	"baliance.com/gooxml/schema/soo/wml"
	"fmt"
)

type Doc struct {
	Driver *document.Document
	With   measurement.Distance
	Height measurement.Distance
	imgRef common.ImageRef
}

func NewDoc() *Doc {
	return &Doc{
		Driver: document.New(),
	}
}

func (self *Doc) SetImgWith(w float64) {
	self.With = measurement.Distance(w * measurement.Centimeter)
	self.Height = self.imgRef.RelativeHeight(self.With)
	para := self.Driver.AddParagraph()
	anchored, err := para.AddRun().AddDrawingAnchored(self.imgRef)
	if err != nil {
		fmt.Printf("unable to add anchored image: %s", err)
	}

	pageW := measurement.Distance(19 * measurement.Centimeter)
	if self.With != 0 {
		pageW = self.With
	}
	pageH := self.imgRef.RelativeHeight(pageW)
	if self.Height != 0 {
		pageH = self.Height
	}
	anchored.SetSize(pageW, pageH)
	anchored.SetOrigin(wml.WdST_RelFromHPage, wml.WdST_RelFromVTopMargin)
	anchored.SetHAlignment(wml.WdST_AlignHCenter)
	anchored.SetYOffset(1 * measurement.Centimeter)
}

func (self *Doc) AddImageToWord(imagePath string) {
	img, err := common.ImageFromFile(imagePath)
	if err != nil {
		fmt.Printf("unable to create image: %s", err)
	}
	iref, err := self.Driver.AddImage(img)
	if err != nil {
		fmt.Printf("unable to add image to document: %s", err)
	}
	self.imgRef = iref
}

func (self *Doc) Save(path string) bool {
	err := self.Driver.SaveToFile(path)
	if err != nil {
		fmt.Printf("unable to save docx: %s", err)
		return false
	}
	return true
}
