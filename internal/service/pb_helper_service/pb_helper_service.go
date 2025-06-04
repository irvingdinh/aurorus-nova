package pb_helper_service

import (
	"github.com/pocketbase/pocketbase/core"
)

type PbHelperService interface {
	EnsureThumbnails(e *core.RecordEvent, fields ...string) error
}

func NewPbHelperService() PbHelperService {
	return &pbHelperService{}
}

type pbHelperService struct {
	//
}

func (s *pbHelperService) EnsureThumbnails(e *core.RecordEvent, fields ...string) error {
	baseLog := e.App.Logger().
		With("func", "EnsureThumbnails").
		With("collection", e.Record.Collection().Name).
		With("id", e.Record.Id)

	fsys, err := e.App.NewFilesystem()
	if err != nil {
		baseLog.With("error", err).Error("Failed to open filesystem")
		return err
	}
	defer fsys.Close()

	for _, fieldName := range fields {
		fieldLog := baseLog.With("field", fieldName)

		field, ok := e.Record.Collection().Fields.GetByName(fieldName).(*core.FileField)
		if !ok {
			fieldLog.Warn("Invalid field")
			continue
		}
		if len(field.Thumbs) == 0 {
			fieldLog.Warn("Invalid thumbs")
			continue
		}

		filenames := e.Record.GetStringSlice(fieldName)
		if len(filenames) == 0 {
			if s := e.Record.GetString(fieldName); s != "" {
				filenames = []string{s}
			}
		}
		if len(filenames) == 0 {
			fieldLog.Debug("No file uploaded, skip the generation")
			continue
		}

		base := e.Record.BaseFilesPath()
		for _, name := range filenames {
			fileLog := fieldLog.With("filename", name)

			original := base + "/" + name
			for _, size := range field.Thumbs {
				sizeLog := fileLog.With("size", size)

				thumb := base + "/thumbs_" + name + "/" + size + "_" + name

				exists, err := fsys.Exists(thumb)
				if err != nil {
					sizeLog.With("error", err).Error("Failed to check existence")
					continue
				}
				if exists {
					sizeLog.Debug("Thumb already exists")
					continue
				}

				if err := fsys.CreateThumb(original, thumb, size); err != nil {
					sizeLog.With("error", err).Error("Failed to create thumb")
				}

				sizeLog.Debug("Created thumbnail successfully")
			}
		}
	}

	return nil
}
