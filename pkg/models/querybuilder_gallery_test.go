// +build integration

package models_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stashapp/stash/pkg/models"
)

func TestGalleryFind(t *testing.T) {
	gqb := models.NewGalleryQueryBuilder()

	const galleryIdx = 0
	gallery, err := gqb.Find(galleryIDs[galleryIdx], nil)

	if err != nil {
		t.Fatalf("Error finding gallery: %s", err.Error())
	}

	assert.Equal(t, getGalleryStringValue(galleryIdx, "Path"), gallery.Path.String)

	gallery, err = gqb.Find(0, nil)

	if err != nil {
		t.Fatalf("Error finding gallery: %s", err.Error())
	}

	assert.Nil(t, gallery)
}

func TestGalleryFindByChecksum(t *testing.T) {
	gqb := models.NewGalleryQueryBuilder()

	const galleryIdx = 0
	galleryChecksum := getGalleryStringValue(galleryIdx, "Checksum")
	gallery, err := gqb.FindByChecksum(galleryChecksum, nil)

	if err != nil {
		t.Fatalf("Error finding gallery: %s", err.Error())
	}

	assert.Equal(t, getGalleryStringValue(galleryIdx, "Path"), gallery.Path.String)

	galleryChecksum = "not exist"
	gallery, err = gqb.FindByChecksum(galleryChecksum, nil)

	if err != nil {
		t.Fatalf("Error finding gallery: %s", err.Error())
	}

	assert.Nil(t, gallery)
}

func TestGalleryFindByPath(t *testing.T) {
	gqb := models.NewGalleryQueryBuilder()

	const galleryIdx = 0
	galleryPath := getGalleryStringValue(galleryIdx, "Path")
	gallery, err := gqb.FindByPath(galleryPath)

	if err != nil {
		t.Fatalf("Error finding gallery: %s", err.Error())
	}

	assert.Equal(t, galleryPath, gallery.Path.String)

	galleryPath = "not exist"
	gallery, err = gqb.FindByPath(galleryPath)

	if err != nil {
		t.Fatalf("Error finding gallery: %s", err.Error())
	}

	assert.Nil(t, gallery)
}

func TestGalleryFindBySceneID(t *testing.T) {
	gqb := models.NewGalleryQueryBuilder()

	sceneID := sceneIDs[sceneIdxWithGallery]
	gallery, err := gqb.FindBySceneID(sceneID, nil)

	if err != nil {
		t.Fatalf("Error finding gallery: %s", err.Error())
	}

	assert.Equal(t, getGalleryStringValue(galleryIdxWithScene, "Path"), gallery.Path.String)

	gallery, err = gqb.FindBySceneID(0, nil)

	if err != nil {
		t.Fatalf("Error finding gallery: %s", err.Error())
	}

	assert.Nil(t, gallery)
}

func TestGalleryQueryQ(t *testing.T) {
	const galleryIdx = 0

	q := getGalleryStringValue(galleryIdx, pathField)

	sqb := models.NewGalleryQueryBuilder()

	galleryQueryQ(t, sqb, q, galleryIdx)
}

func galleryQueryQ(t *testing.T, qb models.GalleryQueryBuilder, q string, expectedGalleryIdx int) {
	filter := models.FindFilterType{
		Q: &q,
	}
	galleries, _ := qb.Query(nil, &filter)

	assert.Len(t, galleries, 1)
	gallery := galleries[0]
	assert.Equal(t, galleryIDs[expectedGalleryIdx], gallery.ID)

	// no Q should return all results
	filter.Q = nil
	galleries, _ = qb.Query(nil, &filter)

	assert.Len(t, galleries, totalGalleries)
}

func TestGalleryQueryPath(t *testing.T) {
	const galleryIdx = 1
	galleryPath := getGalleryStringValue(galleryIdx, "Path")

	pathCriterion := models.StringCriterionInput{
		Value:    galleryPath,
		Modifier: models.CriterionModifierEquals,
	}

	verifyGalleriesPath(t, pathCriterion)

	pathCriterion.Modifier = models.CriterionModifierNotEquals
	verifyGalleriesPath(t, pathCriterion)
}

func verifyGalleriesPath(t *testing.T, pathCriterion models.StringCriterionInput) {
	sqb := models.NewGalleryQueryBuilder()
	galleryFilter := models.GalleryFilterType{
		Path: &pathCriterion,
	}

	galleries, _ := sqb.Query(&galleryFilter, nil)

	for _, gallery := range galleries {
		verifyNullString(t, gallery.Path, pathCriterion)
	}
}

func TestGalleryQueryRating(t *testing.T) {
	const rating = 3
	ratingCriterion := models.IntCriterionInput{
		Value:    rating,
		Modifier: models.CriterionModifierEquals,
	}

	verifyGalleriesRating(t, ratingCriterion)

	ratingCriterion.Modifier = models.CriterionModifierNotEquals
	verifyGalleriesRating(t, ratingCriterion)

	ratingCriterion.Modifier = models.CriterionModifierGreaterThan
	verifyGalleriesRating(t, ratingCriterion)

	ratingCriterion.Modifier = models.CriterionModifierLessThan
	verifyGalleriesRating(t, ratingCriterion)

	ratingCriterion.Modifier = models.CriterionModifierIsNull
	verifyGalleriesRating(t, ratingCriterion)

	ratingCriterion.Modifier = models.CriterionModifierNotNull
	verifyGalleriesRating(t, ratingCriterion)
}

func verifyGalleriesRating(t *testing.T, ratingCriterion models.IntCriterionInput) {
	sqb := models.NewGalleryQueryBuilder()
	galleryFilter := models.GalleryFilterType{
		Rating: &ratingCriterion,
	}

	galleries, _ := sqb.Query(&galleryFilter, nil)

	for _, gallery := range galleries {
		verifyInt64(t, gallery.Rating, ratingCriterion)
	}
}

func TestGalleryQueryIsMissingScene(t *testing.T) {
	qb := models.NewGalleryQueryBuilder()
	isMissing := "scene"
	galleryFilter := models.GalleryFilterType{
		IsMissing: &isMissing,
	}

	q := getGalleryStringValue(galleryIdxWithScene, titleField)
	findFilter := models.FindFilterType{
		Q: &q,
	}

	galleries, _ := qb.Query(&galleryFilter, &findFilter)

	assert.Len(t, galleries, 0)

	findFilter.Q = nil
	galleries, _ = qb.Query(&galleryFilter, &findFilter)

	// ensure non of the ids equal the one with gallery
	for _, gallery := range galleries {
		assert.NotEqual(t, galleryIDs[galleryIdxWithScene], gallery.ID)
	}
}

// TODO ValidGalleriesForScenePath
// TODO Count
// TODO All
// TODO Query
// TODO Update
// TODO Destroy
// TODO ClearGalleryId
