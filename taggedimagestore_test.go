package taggedimages

import (
	"github.com/highlyunavailable/taggedimages/model"
	"github.com/highlyunavailable/taggedimages/store/mem"
	"os"
	"testing"
)

var b *ImageStore

func TestMain(m *testing.M) {
	path := `temp`
	b = New(mem.NewMemDataStore(path))
	os.Exit(m.Run())
}

func Test_ImagePut(t *testing.T) {
	// Costs more than catbus
	f, err := os.Open("testdata/kitty.jpg")
	defer f.Close()
	if err == nil {
		b.PutImage(f, []string{"cat", "monorail"})
	} else {
		t.Fatal(err)
	}

	// A smaller, floatier cat
	f2, err := os.Open("testdata/hover.jpg")
	defer f2.Close()
	if err == nil {
		b.PutImage(f2, []string{"cat", "hover", "kitten"})
	} else {
		t.Fatal(err)
	}

	// Not a cat
	f3, err := os.Open("testdata/bird.jpg")
	defer f3.Close()
	if err == nil {
		b.PutImage(f3, []string{"bird", "bluebird", "grumpy"})
	} else {
		t.Fatal(err)
	}

	// PNG file
	f4, err := os.Open("testdata/lizard.png")
	defer f4.Close()
	if err == nil {
		b.PutImage(f4, []string{"LIZARD", "disappointed"})
	} else {
		t.Fatal(err)
	}

	// Grumpy Cat
	f5, err := os.Open("testdata/grumpycat.jpg")
	defer f5.Close()
	if err == nil {
		b.PutImage(f5, []string{"cat", "grumpy", "painting", "wallpaper"})
	} else {
		t.Fatal(err)
	}
}

func Test_ImageGet(t *testing.T) {
	id := "b2811481cbbf0329c258e042c8f94c0f6509c8c6c42c220537237d105aeabcda"
	img, err := b.GetImage(id)
	if err != nil {
		t.Fatal(err)
	}
	if img.Id != id {
		t.Fatal("Image Id did not match")
	}
}

func Test_ImageList(t *testing.T) {
	imgsP1 := b.GetImages(nil, 1, 3)
	if len(imgsP1) != 3 {
		t.Fatal("Expected 3 images on page 1")
	}

	imgsP2 := b.GetImages(nil, 2, 3)
	if len(imgsP2) != 2 {
		t.Fatal("Expected 2 images on page 2")
	}
}

func Test_ImageFilter(t *testing.T) {
	imgsCat := b.GetImages([]string{"cat"}, 1, 10)
	if len(imgsCat) != 3 {
		t.Fatal("Expected 3 cats")
	}

	imgsGrumpy := b.GetImages([]string{"-grumpy"}, 1, 10)
	if len(imgsGrumpy) != 3 {
		t.Fatal("Expected 3 non-grumpy")
	}

	imgsCatNoGrumpy := b.GetImages([]string{"cat", "-grumpy"}, 1, 10)
	if len(imgsCatNoGrumpy) != 2 {
		t.Fatal("Expected 2 cats that are not grumpy")
	}
}

func Test_ImageTagCasing(t *testing.T) {
	imgsLizard := b.GetImages([]string{"LiZaRd"}, 1, 10)
	if len(imgsLizard) != 1 {
		t.Fatal("Expected 1 lizard")
	}
}

func Test_ImageFilterPaging(t *testing.T) {
	imgsCatP1 := b.GetImages([]string{"cat"}, 1, 2)
	if len(imgsCatP1) != 2 {
		t.Fatal("Expected 2 cats on page 1", len(imgsCatP1))
	}
	imgsCatP2 := b.GetImages([]string{"cat"}, 2, 2)
	if len(imgsCatP2) != 1 {
		t.Fatal("Expected 1 cat on page 2")
	}
}

func Test_ImageDelete(t *testing.T) {
	id := "b2811481cbbf0329c258e042c8f94c0f6509c8c6c42c220537237d105aeabcda"

	err := b.DeleteImage(id)
	if err != nil {
		t.Fatal(err)
	}
	imgsCat := b.GetImages([]string{"cat"}, 1, 10)
	if len(imgsCat) != 2 {
		t.Fatal("Expected 2 cats after deleting 1")
	}
	err = b.DeleteImage(id)
	if err != model.ErrImageNotFound {
		t.Fatal("Expected ErrImageNotFound, got", err)
	}
}
