package tagutils

import (
	"testing"

	"go.aporeto.io/gaia"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTags_TagValidation(t *testing.T) {

	Convey("Given I have a good tag", t, func() {

		tag := gaia.NewTag()
		tag.Value = "alexandre=kind"

		Convey("When I validate it", func() {

			errs := tag.Validate()

			Convey("Then the errs should be nil", func() {
				So(errs, ShouldBeNil)
			})

		})
	})

	Convey("Given I have a good tag with spaces", t, func() {

		tag := gaia.NewTag()
		tag.Value = "alexandre=kind 2"

		Convey("When I validate it", func() {

			errs := tag.Validate()

			Convey("Then the errs should be nil", func() {
				So(errs, ShouldBeNil)
			})

		})
	})

	Convey("Given I have a good tag with spaces and equals", t, func() {

		tag := gaia.NewTag()
		tag.Value = "alexandre=kind 2 = dd"

		Convey("When I validate it", func() {

			errs := tag.Validate()

			Convey("Then the errs should be nil", func() {
				So(errs, ShouldBeNil)
			})

		})
	})

	Convey("Given I have a good tag with no value", t, func() {

		tag := gaia.NewTag()
		tag.Value = ""

		Convey("When I validate it", func() {

			errs := tag.Validate()

			Convey("Then the errs should not be nil", func() {
				So(errs, ShouldNotBeNil)
			})
		})
	})

	Convey("Given I have a tag with a bad value 'key value'", t, func() {

		tag := gaia.NewTag()
		tag.Value = "key value"

		Convey("When I validate it", func() {

			errs := tag.Validate()

			Convey("Then the errs should not be nil", func() {
				So(errs, ShouldNotBeNil)
			})
		})
	})

	Convey("Given I have a tag with a bad value 'key'", t, func() {

		tag := gaia.NewTag()
		tag.Value = "key"

		Convey("When I validate it", func() {

			errs := tag.Validate()

			Convey("Then the errs should not be nil", func() {
				So(errs, ShouldNotBeNil)
			})
		})
	})

	Convey("Given I have a tag with a bad value 'key =value'", t, func() {

		tag := gaia.NewTag()
		tag.Value = "key =value"

		Convey("When I validate it", func() {

			errs := tag.Validate()

			Convey("Then the errs should not be nil", func() {
				So(errs, ShouldNotBeNil)
			})
		})
	})

	Convey("Given I have a tag with a bad value 'key = value'", t, func() {

		tag := gaia.NewTag()
		tag.Value = "key = value"

		Convey("When I validate it", func() {

			errs := tag.Validate()

			Convey("Then the errs should not be nil", func() {
				So(errs, ShouldNotBeNil)
			})
		})
	})

	Convey("Given I have a tag with a bad value 'key key=value'", t, func() {

		tag := gaia.NewTag()
		tag.Value = "key key=value"

		Convey("When I validate it", func() {

			errs := tag.Validate()

			Convey("Then the errs should not be nil", func() {
				So(errs, ShouldNotBeNil)
			})
		})
	})

	Convey("Given I have a tag with a bad value ' key=value'", t, func() {

		tag := gaia.NewTag()
		tag.Value = " key=value"

		Convey("When I validate it", func() {

			errs := tag.Validate()

			Convey("Then the errs should not be nil", func() {
				So(errs, ShouldNotBeNil)
			})
		})
	})

	Convey("Given I have a tag with a bad value ' key=value '", t, func() {

		tag := gaia.NewTag()
		tag.Value = " key=value "

		Convey("When I validate it", func() {

			errs := tag.Validate()

			Convey("Then the errs should not be nil", func() {
				So(errs, ShouldNotBeNil)
			})
		})
	})
}

func TestTags_ValidateTagString(t *testing.T) {

	Convey("Given I have a good tag string ", t, func() {

		str := "alexandre=kind"

		Convey("When I validate it", func() {

			errs := ValidateTagStrings(false, str)

			Convey("Then the errs should be nil", func() {
				So(errs, ShouldBeNil)
			})

		})
	})

	Convey("Given I have an empty tag string", t, func() {

		str := ""

		Convey("When I validate it", func() {

			errs := ValidateTagStrings(false, str)

			Convey("Then the errs should not be nil", func() {
				So(errs, ShouldNotBeNil)
			})
		})
	})

	Convey("Given a tag string 'key value'", t, func() {

		str := "key value"

		Convey("When I validate it", func() {

			errs := ValidateTagStrings(false, str)

			Convey("Then the errs should not be nil", func() {
				So(errs, ShouldNotBeNil)
			})
		})
	})

	Convey("Given a tag string 'key'", t, func() {

		str := "key"

		Convey("When I validate it", func() {

			errs := ValidateTagStrings(false, str)

			Convey("Then the errs should not be nil", func() {
				So(errs, ShouldNotBeNil)
			})
		})
	})

	Convey("Given a tag string 'key =value'", t, func() {

		str := "key =value"

		Convey("When I validate it", func() {

			errs := ValidateTagStrings(false, str)

			Convey("Then the errs should not be nil", func() {
				So(errs, ShouldNotBeNil)
			})
		})
	})

	Convey("Given a tag string 'key = value'", t, func() {

		str := "key = value"

		Convey("When I validate it", func() {

			errs := ValidateTagStrings(false, str)

			Convey("Then the errs should not be nil", func() {
				So(errs, ShouldNotBeNil)
			})
		})
	})

	Convey("Given a tag string 'key key=value'", t, func() {

		str := "key key=value"

		Convey("When I validate it", func() {

			errs := ValidateTagStrings(false, str)

			Convey("Then the errs should not be nil", func() {
				So(errs, ShouldNotBeNil)
			})
		})
	})

	Convey("Given a tag string ' key=value'", t, func() {

		str := " key=value"

		Convey("When I validate it", func() {

			errs := ValidateTagStrings(false, str)

			Convey("Then the errs should not be nil", func() {
				So(errs, ShouldNotBeNil)
			})
		})
	})

	Convey("Given a tag string ' key=value '", t, func() {

		str := " key=value "

		Convey("When I validate it", func() {

			errs := ValidateTagStrings(false, str)

			Convey("Then the errs should not be nil", func() {
				So(errs, ShouldNotBeNil)
			})
		})
	})

	Convey("Given a tag string '@key=value'", t, func() {

		str := "@key=value"

		Convey("When I validate it by allowing system tags", func() {

			errs := ValidateTagStrings(true, str)

			Convey("Then the errs should be nil", func() {
				So(errs, ShouldBeNil)
			})
		})

		Convey("When I validate it without allowing system tags", func() {

			errs := ValidateTagStrings(false, str)

			Convey("Then the errs should not be nil", func() {
				So(errs, ShouldNotBeNil)
			})
		})
	})
}

func TestTags_ValidateMetadataString(t *testing.T) {

	Convey("Given I have a good metadata string ", t, func() {

		str := "@alexandre=kind"

		Convey("When I validate it", func() {

			errs := ValidateMetadataStrings(str)

			Convey("Then the errs should be nil", func() {
				So(errs, ShouldBeNil)
			})

		})
	})

	Convey("Given I have a metadata string not starting with an @", t, func() {

		str := "alexandre=kind"

		Convey("When I validate it", func() {

			errs := ValidateMetadataStrings(str)

			Convey("Then the errs should not be nil", func() {
				So(errs, ShouldNotBeNil)
			})
		})
	})

	Convey("Given I have a metadata string starting with an @auth:", t, func() {

		str := "@auth:alexandre=kind"

		Convey("When I validate it", func() {

			errs := ValidateMetadataStrings(str)

			Convey("Then the errs should not be nil", func() {
				So(errs, ShouldNotBeNil)
			})
		})
	})

	Convey("Given I have a bad tag", t, func() {

		str := "@alexandre kind"

		Convey("When I validate it", func() {

			errs := ValidateMetadataStrings(str)

			Convey("Then the errs should not be nil", func() {
				So(errs, ShouldNotBeNil)
			})
		})
	})
}
