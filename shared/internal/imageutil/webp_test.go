package imageutil_test

import (
	"os"
	"testing"

	"github.com/brunojet/store-go/shared/internal/imageutil"
)

func TestResizeImageToWebP(t *testing.T) {
	inputPath := "D:\\bruno\\OneDrive\\Pictures\\Upscale\\Bem vindos ao RJ vice 16 quase_upscayl_4x_high-fidelity-4x - 98.png"
	imgBytes, err := os.ReadFile(inputPath)
	if err != nil {
		t.Fatalf("Erro ao ler imagem de teste: %v", err)
	}

	tamanhos := []struct {
		largura int
		nome    string
	}{
		{320, "small"},
		{620, "medium"},
		{1024, "large"},
	}

	// Decodifica a imagem original para obter altura
	img, err := imageutil.DecodeImage(imgBytes)
	if err != nil {
		t.Fatalf("Erro ao decodificar imagem: %v", err)
	}
	origW := img.Bounds().Dx()
	origH := img.Bounds().Dy()

	for _, sz := range tamanhos {
		// Calcula altura proporcional
		altura := int(float64(sz.largura) * float64(origH) / float64(origW))
		out, err := imageutil.ResizeImageToPNG(img, sz.largura, altura)
		if err != nil {
			t.Errorf("Erro ao redimensionar para %s: %v", sz.nome, err)
			continue
		}
		outDir := "tmp"
		absPath, err := os.Getwd()
		if err != nil {
			t.Errorf("Erro ao obter diretório atual: %v", err)
		} else {
			t.Logf("Caminho absoluto: %s", absPath+"/"+outDir+"/test_image_"+sz.nome+".png")
		}
		if err := os.MkdirAll(outDir, 0755); err != nil {
			t.Fatalf("Erro ao criar diretório %s: %v", outDir, err)
		}
		outPath := outDir + "/test_image_" + sz.nome + ".png"
		err = os.WriteFile(outPath, out, 0644)
		if err != nil {
			t.Errorf("Erro ao salvar %s: %v", outPath, err)
		}
	}
}
