package models

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/database"
)

type JoinsQueryBuilder struct{}

func NewJoinsQueryBuilder() JoinsQueryBuilder {
	return JoinsQueryBuilder{}
}

func (qb *JoinsQueryBuilder) GetScenePerformers(sceneID int, tx *sqlx.Tx) ([]PerformersScenes, error) {
	ensureTx(tx)

	// Delete the existing joins and then create new ones
	query := `SELECT * from performers_scenes WHERE scene_id = ?`

	var rows *sqlx.Rows
	var err error
	if tx != nil {
		rows, err = tx.Queryx(query, sceneID)
	} else {
		rows, err = database.DB.Queryx(query, sceneID)
	}

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	performerScenes := make([]PerformersScenes, 0)
	for rows.Next() {
		performerScene := PerformersScenes{}
		if err := rows.StructScan(&performerScene); err != nil {
			return nil, err
		}
		performerScenes = append(performerScenes, performerScene)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return performerScenes, nil
}

func (qb *JoinsQueryBuilder) CreatePerformersScenes(newJoins []PerformersScenes, tx *sqlx.Tx) error {
	ensureTx(tx)
	for _, join := range newJoins {
		_, err := tx.NamedExec(
			`INSERT INTO performers_scenes (performer_id, scene_id) VALUES (:performer_id, :scene_id)`,
			join,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// AddPerformerScene adds a performer to a scene. It does not make any change
// if the performer already exists on the scene. It returns true if scene
// performer was added.
func (qb *JoinsQueryBuilder) AddPerformerScene(sceneID int, performerID int, tx *sqlx.Tx) (bool, error) {
	ensureTx(tx)

	existingPerformers, err := qb.GetScenePerformers(sceneID, tx)

	if err != nil {
		return false, err
	}

	// ensure not already present
	for _, p := range existingPerformers {
		if p.PerformerID == performerID && p.SceneID == sceneID {
			return false, nil
		}
	}

	performerJoin := PerformersScenes{
		PerformerID: performerID,
		SceneID:     sceneID,
	}
	performerJoins := append(existingPerformers, performerJoin)

	err = qb.UpdatePerformersScenes(sceneID, performerJoins, tx)

	return err == nil, err
}

func (qb *JoinsQueryBuilder) UpdatePerformersScenes(sceneID int, updatedJoins []PerformersScenes, tx *sqlx.Tx) error {
	ensureTx(tx)

	// Delete the existing joins and then create new ones
	_, err := tx.Exec("DELETE FROM performers_scenes WHERE scene_id = ?", sceneID)
	if err != nil {
		return err
	}
	return qb.CreatePerformersScenes(updatedJoins, tx)
}

func (qb *JoinsQueryBuilder) DestroyPerformersScenes(sceneID int, tx *sqlx.Tx) error {
	ensureTx(tx)

	// Delete the existing joins
	_, err := tx.Exec("DELETE FROM performers_scenes WHERE scene_id = ?", sceneID)
	return err
}

func (qb *JoinsQueryBuilder) GetSceneMovies(sceneID int, tx *sqlx.Tx) ([]MoviesScenes, error) {
	query := `SELECT * from movies_scenes WHERE scene_id = ?`

	var rows *sqlx.Rows
	var err error
	if tx != nil {
		rows, err = tx.Queryx(query, sceneID)
	} else {
		rows, err = database.DB.Queryx(query, sceneID)
	}

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	movieScenes := make([]MoviesScenes, 0)
	for rows.Next() {
		movieScene := MoviesScenes{}
		if err := rows.StructScan(&movieScene); err != nil {
			return nil, err
		}
		movieScenes = append(movieScenes, movieScene)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return movieScenes, nil
}

func (qb *JoinsQueryBuilder) CreateMoviesScenes(newJoins []MoviesScenes, tx *sqlx.Tx) error {
	ensureTx(tx)
	for _, join := range newJoins {
		_, err := tx.NamedExec(
			`INSERT INTO movies_scenes (movie_id, scene_id, scene_index) VALUES (:movie_id, :scene_id, :scene_index)`,
			join,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// AddMovieScene adds a movie to a scene. It does not make any change
// if the movie already exists on the scene. It returns true if scene
// movie was added.

func (qb *JoinsQueryBuilder) AddMoviesScene(sceneID int, movieID int, sceneIdx *int, tx *sqlx.Tx) (bool, error) {
	ensureTx(tx)

	existingMovies, err := qb.GetSceneMovies(sceneID, tx)

	if err != nil {
		return false, err
	}

	// ensure not already present
	for _, p := range existingMovies {
		if p.MovieID == movieID && p.SceneID == sceneID {
			return false, nil
		}
	}

	movieJoin := MoviesScenes{
		MovieID: movieID,
		SceneID: sceneID,
	}

	if sceneIdx != nil {
		movieJoin.SceneIndex = sql.NullInt64{
			Int64: int64(*sceneIdx),
			Valid: true,
		}
	}
	movieJoins := append(existingMovies, movieJoin)

	err = qb.UpdateMoviesScenes(sceneID, movieJoins, tx)

	return err == nil, err
}

func (qb *JoinsQueryBuilder) UpdateMoviesScenes(sceneID int, updatedJoins []MoviesScenes, tx *sqlx.Tx) error {
	ensureTx(tx)

	// Delete the existing joins and then create new ones
	_, err := tx.Exec("DELETE FROM movies_scenes WHERE scene_id = ?", sceneID)
	if err != nil {
		return err
	}
	return qb.CreateMoviesScenes(updatedJoins, tx)
}

func (qb *JoinsQueryBuilder) DestroyMoviesScenes(sceneID int, tx *sqlx.Tx) error {
	ensureTx(tx)

	// Delete the existing joins
	_, err := tx.Exec("DELETE FROM movies_scenes WHERE scene_id = ?", sceneID)
	return err
}

func (qb *JoinsQueryBuilder) GetSceneTags(sceneID int, tx *sqlx.Tx) ([]ScenesTags, error) {
	ensureTx(tx)

	// Delete the existing joins and then create new ones
	query := `SELECT * from scenes_tags WHERE scene_id = ?`

	var rows *sqlx.Rows
	var err error
	if tx != nil {
		rows, err = tx.Queryx(query, sceneID)
	} else {
		rows, err = database.DB.Queryx(query, sceneID)
	}

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	sceneTags := make([]ScenesTags, 0)
	for rows.Next() {
		sceneTag := ScenesTags{}
		if err := rows.StructScan(&sceneTag); err != nil {
			return nil, err
		}
		sceneTags = append(sceneTags, sceneTag)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sceneTags, nil
}

func (qb *JoinsQueryBuilder) CreateScenesTags(newJoins []ScenesTags, tx *sqlx.Tx) error {
	ensureTx(tx)
	for _, join := range newJoins {
		_, err := tx.NamedExec(
			`INSERT INTO scenes_tags (scene_id, tag_id) VALUES (:scene_id, :tag_id)`,
			join,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (qb *JoinsQueryBuilder) UpdateScenesTags(sceneID int, updatedJoins []ScenesTags, tx *sqlx.Tx) error {
	ensureTx(tx)

	// Delete the existing joins and then create new ones
	_, err := tx.Exec("DELETE FROM scenes_tags WHERE scene_id = ?", sceneID)
	if err != nil {
		return err
	}
	return qb.CreateScenesTags(updatedJoins, tx)
}

// AddSceneTag adds a tag to a scene. It does not make any change if the tag
// already exists on the scene. It returns true if scene tag was added.
func (qb *JoinsQueryBuilder) AddSceneTag(sceneID int, tagID int, tx *sqlx.Tx) (bool, error) {
	ensureTx(tx)

	existingTags, err := qb.GetSceneTags(sceneID, tx)

	if err != nil {
		return false, err
	}

	// ensure not already present
	for _, p := range existingTags {
		if p.TagID == tagID && p.SceneID == sceneID {
			return false, nil
		}
	}

	tagJoin := ScenesTags{
		TagID:   tagID,
		SceneID: sceneID,
	}
	tagJoins := append(existingTags, tagJoin)

	err = qb.UpdateScenesTags(sceneID, tagJoins, tx)

	return err == nil, err
}

func (qb *JoinsQueryBuilder) DestroyScenesTags(sceneID int, tx *sqlx.Tx) error {
	ensureTx(tx)

	// Delete the existing joins
	_, err := tx.Exec("DELETE FROM scenes_tags WHERE scene_id = ?", sceneID)

	return err
}

func (qb *JoinsQueryBuilder) CreateSceneMarkersTags(newJoins []SceneMarkersTags, tx *sqlx.Tx) error {
	ensureTx(tx)
	for _, join := range newJoins {
		_, err := tx.NamedExec(
			`INSERT INTO scene_markers_tags (scene_marker_id, tag_id) VALUES (:scene_marker_id, :tag_id)`,
			join,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (qb *JoinsQueryBuilder) UpdateSceneMarkersTags(sceneMarkerID int, updatedJoins []SceneMarkersTags, tx *sqlx.Tx) error {
	ensureTx(tx)

	// Delete the existing joins and then create new ones
	_, err := tx.Exec("DELETE FROM scene_markers_tags WHERE scene_marker_id = ?", sceneMarkerID)
	if err != nil {
		return err
	}
	return qb.CreateSceneMarkersTags(updatedJoins, tx)
}

func (qb *JoinsQueryBuilder) DestroySceneMarkersTags(sceneMarkerID int, updatedJoins []SceneMarkersTags, tx *sqlx.Tx) error {
	ensureTx(tx)

	// Delete the existing joins
	_, err := tx.Exec("DELETE FROM scene_markers_tags WHERE scene_marker_id = ?", sceneMarkerID)
	return err
}

func (qb *JoinsQueryBuilder) DestroyScenesGalleries(sceneID int, tx *sqlx.Tx) error {
	ensureTx(tx)

	// Unset the existing scene id from galleries
	_, err := tx.Exec("UPDATE galleries SET scene_id = null WHERE scene_id = ?", sceneID)

	return err
}

func (qb *JoinsQueryBuilder) DestroyScenesMarkers(sceneID int, tx *sqlx.Tx) error {
	ensureTx(tx)

	// Delete the scene marker tags
	_, err := tx.Exec("DELETE t FROM scene_markers_tags t join scene_markers m on t.scene_marker_id = m.id WHERE m.scene_id = ?", sceneID)

	// Delete the existing joins
	_, err = tx.Exec("DELETE FROM scene_markers WHERE scene_id = ?", sceneID)

	return err
}

func (qb *JoinsQueryBuilder) CreateStashIDs(entityName string, entityID int, newJoins []StashID, tx *sqlx.Tx) error {
	query := "INSERT INTO " + entityName + "_stash_ids (" + entityName + "_id, endpoint, stash_id) VALUES (?, ?, ?)"
	ensureTx(tx)
	for _, join := range newJoins {
		_, err := tx.Exec(query, entityID, join.Endpoint, join.StashID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (qb *JoinsQueryBuilder) GetImagePerformers(imageID int, tx *sqlx.Tx) ([]PerformersImages, error) {
	ensureTx(tx)

	// Delete the existing joins and then create new ones
	query := `SELECT * from performers_images WHERE image_id = ?`

	var rows *sqlx.Rows
	var err error
	if tx != nil {
		rows, err = tx.Queryx(query, imageID)
	} else {
		rows, err = database.DB.Queryx(query, imageID)
	}

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	performerImages := make([]PerformersImages, 0)
	for rows.Next() {
		performerImage := PerformersImages{}
		if err := rows.StructScan(&performerImage); err != nil {
			return nil, err
		}
		performerImages = append(performerImages, performerImage)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return performerImages, nil
}

func (qb *JoinsQueryBuilder) CreatePerformersImages(newJoins []PerformersImages, tx *sqlx.Tx) error {
	ensureTx(tx)
	for _, join := range newJoins {
		_, err := tx.NamedExec(
			`INSERT INTO performers_images (performer_id, image_id) VALUES (:performer_id, :image_id)`,
			join,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// AddPerformerImage adds a performer to a image. It does not make any change
// if the performer already exists on the image. It returns true if image
// performer was added.
func (qb *JoinsQueryBuilder) AddPerformerImage(imageID int, performerID int, tx *sqlx.Tx) (bool, error) {
	ensureTx(tx)

	existingPerformers, err := qb.GetImagePerformers(imageID, tx)

	if err != nil {
		return false, err
	}

	// ensure not already present
	for _, p := range existingPerformers {
		if p.PerformerID == performerID && p.ImageID == imageID {
			return false, nil
		}
	}

	performerJoin := PerformersImages{
		PerformerID: performerID,
		ImageID:     imageID,
	}
	performerJoins := append(existingPerformers, performerJoin)

	err = qb.UpdatePerformersImages(imageID, performerJoins, tx)

	return err == nil, err
}

func (qb *JoinsQueryBuilder) UpdatePerformersImages(imageID int, updatedJoins []PerformersImages, tx *sqlx.Tx) error {
	ensureTx(tx)

	// Delete the existing joins and then create new ones
	_, err := tx.Exec("DELETE FROM performers_images WHERE image_id = ?", imageID)
	if err != nil {
		return err
	}
	return qb.CreatePerformersImages(updatedJoins, tx)
}

func (qb *JoinsQueryBuilder) DestroyPerformersImages(imageID int, tx *sqlx.Tx) error {
	ensureTx(tx)

	// Delete the existing joins
	_, err := tx.Exec("DELETE FROM performers_images WHERE image_id = ?", imageID)
	return err
}

func (qb *JoinsQueryBuilder) GetImageTags(imageID int, tx *sqlx.Tx) ([]ImagesTags, error) {
	ensureTx(tx)

	// Delete the existing joins and then create new ones
	query := `SELECT * from images_tags WHERE image_id = ?`

	var rows *sqlx.Rows
	var err error
	if tx != nil {
		rows, err = tx.Queryx(query, imageID)
	} else {
		rows, err = database.DB.Queryx(query, imageID)
	}

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	imageTags := make([]ImagesTags, 0)
	for rows.Next() {
		imageTag := ImagesTags{}
		if err := rows.StructScan(&imageTag); err != nil {
			return nil, err
		}
		imageTags = append(imageTags, imageTag)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return imageTags, nil
}

func (qb *JoinsQueryBuilder) CreateImagesTags(newJoins []ImagesTags, tx *sqlx.Tx) error {
	ensureTx(tx)
	for _, join := range newJoins {
		_, err := tx.NamedExec(
			`INSERT INTO images_tags (image_id, tag_id) VALUES (:image_id, :tag_id)`,
			join,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (qb *JoinsQueryBuilder) UpdateImagesTags(imageID int, updatedJoins []ImagesTags, tx *sqlx.Tx) error {
	ensureTx(tx)

	// Delete the existing joins and then create new ones
	_, err := tx.Exec("DELETE FROM images_tags WHERE image_id = ?", imageID)
	if err != nil {
		return err
	}
	return qb.CreateImagesTags(updatedJoins, tx)
}

// AddImageTag adds a tag to a image. It does not make any change if the tag
// already exists on the image. It returns true if image tag was added.
func (qb *JoinsQueryBuilder) AddImageTag(imageID int, tagID int, tx *sqlx.Tx) (bool, error) {
	ensureTx(tx)

	existingTags, err := qb.GetImageTags(imageID, tx)

	if err != nil {
		return false, err
	}

	// ensure not already present
	for _, p := range existingTags {
		if p.TagID == tagID && p.ImageID == imageID {
			return false, nil
		}
	}

	tagJoin := ImagesTags{
		TagID:   tagID,
		ImageID: imageID,
	}
	tagJoins := append(existingTags, tagJoin)

	err = qb.UpdateImagesTags(imageID, tagJoins, tx)

	return err == nil, err
}

func (qb *JoinsQueryBuilder) DestroyImagesTags(imageID int, tx *sqlx.Tx) error {
	ensureTx(tx)

	// Delete the existing joins
	_, err := tx.Exec("DELETE FROM images_tags WHERE image_id = ?", imageID)

	return err
}

func (qb *JoinsQueryBuilder) GetImageGalleries(imageID int, tx *sqlx.Tx) ([]GalleriesImages, error) {
	ensureTx(tx)

	// Delete the existing joins and then create new ones
	query := `SELECT * from galleries_images WHERE image_id = ?`

	var rows *sqlx.Rows
	var err error
	if tx != nil {
		rows, err = tx.Queryx(query, imageID)
	} else {
		rows, err = database.DB.Queryx(query, imageID)
	}

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	galleryImages := make([]GalleriesImages, 0)
	for rows.Next() {
		galleriesImages := GalleriesImages{}
		if err := rows.StructScan(&galleriesImages); err != nil {
			return nil, err
		}
		galleryImages = append(galleryImages, galleriesImages)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return galleryImages, nil
}

func (qb *JoinsQueryBuilder) CreateGalleriesImages(newJoins []GalleriesImages, tx *sqlx.Tx) error {
	ensureTx(tx)
	for _, join := range newJoins {
		_, err := tx.NamedExec(
			`INSERT INTO galleries_images (gallery_id, image_id) VALUES (:gallery_id, :image_id)`,
			join,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (qb *JoinsQueryBuilder) UpdateGalleriesImages(imageID int, updatedJoins []GalleriesImages, tx *sqlx.Tx) error {
	ensureTx(tx)

	// Delete the existing joins and then create new ones
	_, err := tx.Exec("DELETE FROM galleries_images WHERE image_id = ?", imageID)
	if err != nil {
		return err
	}
	return qb.CreateGalleriesImages(updatedJoins, tx)
}

// AddGalleryImage adds a gallery to an image. It does not make any change if the tag
// already exists on the image. It returns true if image tag was added.
func (qb *JoinsQueryBuilder) AddImageGallery(imageID int, galleryID int, tx *sqlx.Tx) (bool, error) {
	ensureTx(tx)

	existingGalleries, err := qb.GetImageGalleries(imageID, tx)

	if err != nil {
		return false, err
	}

	// ensure not already present
	for _, p := range existingGalleries {
		if p.GalleryID == galleryID && p.ImageID == imageID {
			return false, nil
		}
	}

	galleryJoin := GalleriesImages{
		GalleryID: galleryID,
		ImageID:   imageID,
	}
	galleryJoins := append(existingGalleries, galleryJoin)

	err = qb.UpdateGalleriesImages(imageID, galleryJoins, tx)

	return err == nil, err
}

// RemoveImageGallery removes a gallery from an image. Returns true if the join
// was removed.
func (qb *JoinsQueryBuilder) RemoveImageGallery(imageID int, galleryID int, tx *sqlx.Tx) (bool, error) {
	ensureTx(tx)

	existingGalleries, err := qb.GetImageGalleries(imageID, tx)

	if err != nil {
		return false, err
	}

	// remove the join
	var updatedJoins []GalleriesImages
	found := false
	for _, p := range existingGalleries {
		if p.GalleryID == galleryID && p.ImageID == imageID {
			found = true
			continue
		}

		updatedJoins = append(updatedJoins, p)
	}

	if found {
		err = qb.UpdateGalleriesImages(imageID, updatedJoins, tx)
	}

	return found && err == nil, err
}

func (qb *JoinsQueryBuilder) DestroyImageGalleries(imageID int, tx *sqlx.Tx) error {
	ensureTx(tx)

	// Delete the existing joins
	_, err := tx.Exec("DELETE FROM galleries_images WHERE image_id = ?", imageID)

	return err
}

func (qb *JoinsQueryBuilder) GetGalleryPerformers(galleryID int, tx *sqlx.Tx) ([]PerformersGalleries, error) {
	ensureTx(tx)

	// Delete the existing joins and then create new ones
	query := `SELECT * from performers_galleries WHERE gallery_id = ?`

	var rows *sqlx.Rows
	var err error
	if tx != nil {
		rows, err = tx.Queryx(query, galleryID)
	} else {
		rows, err = database.DB.Queryx(query, galleryID)
	}

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	performerGalleries := make([]PerformersGalleries, 0)
	for rows.Next() {
		performerGallery := PerformersGalleries{}
		if err := rows.StructScan(&performerGallery); err != nil {
			return nil, err
		}
		performerGalleries = append(performerGalleries, performerGallery)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return performerGalleries, nil
}

func (qb *JoinsQueryBuilder) CreatePerformersGalleries(newJoins []PerformersGalleries, tx *sqlx.Tx) error {
	ensureTx(tx)
	for _, join := range newJoins {
		_, err := tx.NamedExec(
			`INSERT INTO performers_galleries (performer_id, gallery_id) VALUES (:performer_id, :gallery_id)`,
			join,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// AddPerformerGallery adds a performer to a gallery. It does not make any change
// if the performer already exists on the gallery. It returns true if gallery
// performer was added.
func (qb *JoinsQueryBuilder) AddPerformerGallery(galleryID int, performerID int, tx *sqlx.Tx) (bool, error) {
	ensureTx(tx)

	existingPerformers, err := qb.GetGalleryPerformers(galleryID, tx)

	if err != nil {
		return false, err
	}

	// ensure not already present
	for _, p := range existingPerformers {
		if p.PerformerID == performerID && p.GalleryID == galleryID {
			return false, nil
		}
	}

	performerJoin := PerformersGalleries{
		PerformerID: performerID,
		GalleryID:   galleryID,
	}
	performerJoins := append(existingPerformers, performerJoin)

	err = qb.UpdatePerformersGalleries(galleryID, performerJoins, tx)

	return err == nil, err
}

func (qb *JoinsQueryBuilder) UpdatePerformersGalleries(galleryID int, updatedJoins []PerformersGalleries, tx *sqlx.Tx) error {
	ensureTx(tx)

	// Delete the existing joins and then create new ones
	_, err := tx.Exec("DELETE FROM performers_galleries WHERE gallery_id = ?", galleryID)
	if err != nil {
		return err
	}
	return qb.CreatePerformersGalleries(updatedJoins, tx)
}

func (qb *JoinsQueryBuilder) DestroyPerformersGalleries(galleryID int, tx *sqlx.Tx) error {
	ensureTx(tx)

	// Delete the existing joins
	_, err := tx.Exec("DELETE FROM performers_galleries WHERE gallery_id = ?", galleryID)
	return err
}

func (qb *JoinsQueryBuilder) GetGalleryTags(galleryID int, tx *sqlx.Tx) ([]GalleriesTags, error) {
	ensureTx(tx)

	// Delete the existing joins and then create new ones
	query := `SELECT * from galleries_tags WHERE gallery_id = ?`

	var rows *sqlx.Rows
	var err error
	if tx != nil {
		rows, err = tx.Queryx(query, galleryID)
	} else {
		rows, err = database.DB.Queryx(query, galleryID)
	}

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	galleryTags := make([]GalleriesTags, 0)
	for rows.Next() {
		galleryTag := GalleriesTags{}
		if err := rows.StructScan(&galleryTag); err != nil {
			return nil, err
		}
		galleryTags = append(galleryTags, galleryTag)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return galleryTags, nil
}

func (qb *JoinsQueryBuilder) CreateGalleriesTags(newJoins []GalleriesTags, tx *sqlx.Tx) error {
	ensureTx(tx)
	for _, join := range newJoins {
		_, err := tx.NamedExec(
			`INSERT INTO galleries_tags (gallery_id, tag_id) VALUES (:gallery_id, :tag_id)`,
			join,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (qb *JoinsQueryBuilder) UpdateGalleriesTags(galleryID int, updatedJoins []GalleriesTags, tx *sqlx.Tx) error {
	ensureTx(tx)

	// Delete the existing joins and then create new ones
	_, err := tx.Exec("DELETE FROM galleries_tags WHERE gallery_id = ?", galleryID)
	if err != nil {
		return err
	}
	return qb.CreateGalleriesTags(updatedJoins, tx)
}

// AddGalleryTag adds a tag to a gallery. It does not make any change if the tag
// already exists on the gallery. It returns true if gallery tag was added.
func (qb *JoinsQueryBuilder) AddGalleryTag(galleryID int, tagID int, tx *sqlx.Tx) (bool, error) {
	ensureTx(tx)

	existingTags, err := qb.GetGalleryTags(galleryID, tx)

	if err != nil {
		return false, err
	}

	// ensure not already present
	for _, p := range existingTags {
		if p.TagID == tagID && p.GalleryID == galleryID {
			return false, nil
		}
	}

	tagJoin := GalleriesTags{
		TagID:     tagID,
		GalleryID: galleryID,
	}
	tagJoins := append(existingTags, tagJoin)

	err = qb.UpdateGalleriesTags(galleryID, tagJoins, tx)

	return err == nil, err
}

func (qb *JoinsQueryBuilder) DestroyGalleriesTags(galleryID int, tx *sqlx.Tx) error {
	ensureTx(tx)

	// Delete the existing joins
	_, err := tx.Exec("DELETE FROM galleries_tags WHERE gallery_id = ?", galleryID)

	return err
}

func (qb *JoinsQueryBuilder) GetSceneStashIDs(sceneID int) ([]*StashID, error) {
	rows, err := database.DB.Queryx(`SELECT stash_id, endpoint from scene_stash_ids WHERE scene_id = ?`, sceneID)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	stashIDs := []*StashID{}
	for rows.Next() {
		stashID := StashID{}
		if err := rows.StructScan(&stashID); err != nil {
			return nil, err
		}
		stashIDs = append(stashIDs, &stashID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return stashIDs, nil
}

func (qb *JoinsQueryBuilder) GetPerformerStashIDs(performerID int) ([]*StashID, error) {
	rows, err := database.DB.Queryx(`SELECT stash_id, endpoint from performer_stash_ids WHERE performer_id = ?`, performerID)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	stashIDs := []*StashID{}
	for rows.Next() {
		stashID := StashID{}
		if err := rows.StructScan(&stashID); err != nil {
			return nil, err
		}
		stashIDs = append(stashIDs, &stashID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return stashIDs, nil
}

func (qb *JoinsQueryBuilder) GetStudioStashIDs(studioID int) ([]*StashID, error) {
	rows, err := database.DB.Queryx(`SELECT stash_id, endpoint from studio_stash_ids WHERE studio_id = ?`, studioID)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	stashIDs := []*StashID{}
	for rows.Next() {
		stashID := StashID{}
		if err := rows.StructScan(&stashID); err != nil {
			return nil, err
		}
		stashIDs = append(stashIDs, &stashID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return stashIDs, nil
}

func (qb *JoinsQueryBuilder) UpdateSceneStashIDs(sceneID int, updatedJoins []StashID, tx *sqlx.Tx) error {
	ensureTx(tx)

	_, err := tx.Exec("DELETE FROM scene_stash_ids WHERE scene_id = ?", sceneID)
	if err != nil {
		return err
	}
	return qb.CreateStashIDs("scene", sceneID, updatedJoins, tx)
}

func (qb *JoinsQueryBuilder) UpdatePerformerStashIDs(performerID int, updatedJoins []StashID, tx *sqlx.Tx) error {
	ensureTx(tx)

	_, err := tx.Exec("DELETE FROM performer_stash_ids WHERE performer_id = ?", performerID)
	if err != nil {
		return err
	}
	return qb.CreateStashIDs("performer", performerID, updatedJoins, tx)
}

func (qb *JoinsQueryBuilder) UpdateStudioStashIDs(studioID int, updatedJoins []StashID, tx *sqlx.Tx) error {
	ensureTx(tx)

	_, err := tx.Exec("DELETE FROM studio_stash_ids WHERE studio_id = ?", studioID)
	if err != nil {
		return err
	}
	return qb.CreateStashIDs("studio", studioID, updatedJoins, tx)
}
