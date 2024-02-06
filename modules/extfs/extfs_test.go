package extfs_test

import (
	"bytes"
	"io/fs"
	"os"
	"pan/modules/extfs"
	"pan/modules/extfs/models"
	"pan/modules/extfs/repositories"
	"pan/modules/extfs/services"
	"pan/peer"
	"path"
	"testing"

	mockedEvent "pan/mocks/pan/modules/extfs/events"
	mockedPeer "pan/mocks/pan/modules/extfs/peer"
	mockedRepo "pan/mocks/pan/modules/extfs/repositories"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestExtFS ...
func TestExtFS(t *testing.T) {

	// setup ...
	setup := func() *extfs.ExtFS {
		efs := new(extfs.ExtFS)
		efs.RemotePeerService = new(services.RemotePeerService)
		efs.RemoteFilesStateService = new(services.RemoteFilesStateService)
		efs.TargetService = new(services.TargetService)
		efs.TargetService.TargetFileService = new(services.TargetFileService)
		return efs
	}

	setupTempFS := func() string {
		dir, err := os.MkdirTemp(os.TempDir(), "extfs-test")
		if err != nil {
			panic(err)
		}
		return dir
	}

	setupTempDir := func(dirPath string) ([]string, []fs.FileInfo, []byte) {

		firstLevelDirPath, err := os.MkdirTemp(dirPath, "first-level")
		if err != nil {
			panic(err)
		}

		secondLevelDirPath, err := os.MkdirTemp(firstLevelDirPath, "second-level")
		if err != nil {
			panic(err)
		}

		fileName := "test.txt"
		fileContent := []byte("ExtFS Test File Content!")
		os.WriteFile(path.Join(dirPath, fileName), fileContent, 0644)
		os.WriteFile(path.Join(firstLevelDirPath, fileName), fileContent, 0644)
		os.WriteFile(path.Join(secondLevelDirPath, fileName), fileContent, 0644)

		dirFS := os.DirFS(dirPath)
		filePaths := make([]string, 0)
		fileStats := make([]fs.FileInfo, 0)
		fs.WalkDir(dirFS, ".", func(path string, d fs.DirEntry, err error) error {
			if err == nil && d.IsDir() == false {
				filePaths = append(filePaths, path)
				stat, err := d.Info()
				if err == nil {
					fileStats = append(fileStats, stat)
				}
			}
			return err
		})

		return filePaths, fileStats, fileContent
	}

	teardownTempFS := func(path string) {
		os.RemoveAll(path)
	}

	t.Run("OnInit", func(t *testing.T) {
		dirPath := setupTempFS()
		defer teardownTempFS(dirPath)
		filePaths, fileStats, fileContent := setupTempDir(dirPath)
		dirStat, err := os.Stat(dirPath)
		if err != nil {
			t.Fatal(err)
		}
		dirSize := int64(0)
		for _, stat := range fileStats {
			dirSize += stat.Size()
		}

		efs := setup()
		sig, err := efs.TargetService.TargetFileService.GenerateFileSignature(bytes.NewReader(fileContent))
		if err != nil {
			t.Fatal(err)
		}

		targetRepo := new(mockedRepo.MockTargetRepository)
		defer targetRepo.AssertExpectations(t)
		efs.TargetService.TargetRepo = targetRepo

		var targetRow models.Target
		targetRow.ID = uint(123)
		targetRow.FilePath = dirPath
		targetRow.Enabled = true
		targetRepo.On("FindAllWithEnabled").Once().Return([]models.Target{targetRow}, nil)
		sTargetRow := targetRow
		sTargetRow.Name = dirStat.Name()
		sTargetRow.ModifyTime = dirStat.ModTime()
		sTargetRow.Total = uint(len(filePaths))
		sTargetRow.Size = dirSize
		targetRepo.On("Save", sTargetRow).Once().Return(nil)

		TargetFileRepo := new(mockedRepo.MockTargetFileRepository)
		defer TargetFileRepo.AssertExpectations(t)
		efs.TargetService.TargetFileService.TargetFileRepo = TargetFileRepo

		TargetFileRepo.On("UpdateEachFileInfoByTargetID", targetRow.ID, mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
			fileInfoIteration := args.Get(1).(repositories.FileInfoIteration)
			var fileInfo models.TargetFile
			fileInfo.RelativePath = filePaths[0]
			err := fileInfoIteration(&fileInfo)

			assert.Nil(t, err)
			assert.Equal(t, fileStats[0].ModTime(), fileInfo.ModifyTime)
			assert.Equal(t, fileStats[0].Size(), fileInfo.Size)
			assert.Equal(t, sig, fileInfo.Hash)
			assert.Equal(t, fileStats[0].Name(), fileInfo.Name)

			var missFileInfo models.TargetFile
			missFileInfo.RelativePath = filePaths[0] + "_miss"
			err = fileInfoIteration(&missFileInfo)

			assert.Equal(t, fs.ErrNotExist, err)

		})

		var firstFileInfo models.TargetFile
		firstFileInfo.TargetID = targetRow.ID
		firstFileInfo.RelativePath = filePaths[0]
		TargetFileRepo.On("FindOrCreateByTargetIDAndRelativePath", targetRow.ID, firstFileInfo.RelativePath).Once().Return(firstFileInfo, nil)
		sFirstFileInfo := firstFileInfo
		sFirstFileInfo.Hash = sig
		sFirstFileInfo.ModifyTime = fileStats[0].ModTime()
		sFirstFileInfo.Size = fileStats[0].Size()
		sFirstFileInfo.Name = fileStats[0].Name()
		TargetFileRepo.On("Save", sFirstFileInfo).Once().Return(nil)
		var secondFileInfo models.TargetFile
		secondFileInfo.TargetID = targetRow.ID
		secondFileInfo.RelativePath = filePaths[1]
		TargetFileRepo.On("FindOrCreateByTargetIDAndRelativePath", targetRow.ID, secondFileInfo.RelativePath).Once().Return(secondFileInfo, nil)
		sSecondFileInfo := secondFileInfo
		sSecondFileInfo.Hash = sig
		sSecondFileInfo.ModifyTime = fileStats[1].ModTime()
		sSecondFileInfo.Size = fileStats[1].Size()
		sSecondFileInfo.Name = fileStats[1].Name()
		TargetFileRepo.On("Save", sSecondFileInfo).Once().Return(nil)
		var thirdFileInfo models.TargetFile
		thirdFileInfo.TargetID = targetRow.ID
		thirdFileInfo.RelativePath = filePaths[2]
		TargetFileRepo.On("FindOrCreateByTargetIDAndRelativePath", targetRow.ID, thirdFileInfo.RelativePath).Once().Return(thirdFileInfo, nil)
		sThirdFileInfo := thirdFileInfo
		sThirdFileInfo.Hash = sig
		sThirdFileInfo.ModifyTime = fileStats[2].ModTime()
		sThirdFileInfo.Size = fileStats[2].Size()
		sThirdFileInfo.Name = fileStats[2].Name()
		TargetFileRepo.On("Save", sThirdFileInfo).Once().Return(nil)

		efs.OnInit()

	})

	t.Run("OnNodeAdded with enabled", func(t *testing.T) {

		peerId := peer.PeerId(uuid.New())
		efs := setup()

		remotePeerRepo := new(mockedRepo.MockRemotePeerRepository)
		defer remotePeerRepo.AssertExpectations(t)
		efs.RemotePeerService.RemotePeerRepo = remotePeerRepo
		var remotePeerRow models.RemotePeer
		remotePeerRow.Enabled = true
		remotePeerRepo.On("FindOne", peerId.String()).Once().Return(remotePeerRow, nil)

		filesStateRepo := new(mockedRepo.MockRemoteFilesStateRepository)
		defer filesStateRepo.AssertExpectations(t)
		efs.RemoteFilesStateService.RemoteFilesStateRepo = filesStateRepo
		var stateRow models.RemoteFilesState
		stateRow.RemoteHash = []byte{1, 2, 3, 4, 5, 6, 7}
		filesStateRepo.On("FindOne", peerId.String()).Once().Return(stateRow, nil)

		api := new(mockedPeer.MockAPI)
		efs.RemoteFilesStateService.API = api
		defer api.AssertExpectations(t)
		var stateInfo models.RemoteStateInfo
		stateInfo.Hash = []byte{11, 12, 13, 14, 15, 16, 17}
		stateInfo.Time = 123
		api.On("GetRemoteFilesState", peerId).Once().Return(stateInfo, nil)

		stateRow.ID = peerId.String()
		stateRow.RemoteHash = stateInfo.Hash
		stateRow.RemoteTime = stateInfo.Time
		filesStateRepo.On("Save", stateRow).Once().Return(nil)

		event := new(mockedEvent.MockRemoteFilesStateEvent)
		defer event.AssertExpectations(t)
		efs.RemoteFilesStateService.RemoteFilesStateEvent = event
		event.On("OnRemoteFilesStateUpdated", peerId).Once()

		efs.OnNodeAdded(peerId)

	})

	t.Run("OnNodeAdded with disabled", func(t *testing.T) {
		peerId := peer.PeerId(uuid.New())
		efs := setup()

		remotePeerRepo := new(mockedRepo.MockRemotePeerRepository)
		defer remotePeerRepo.AssertExpectations(t)
		efs.RemotePeerService.RemotePeerRepo = remotePeerRepo
		var remotePeerRow models.RemotePeer
		remotePeerRow.Enabled = false
		remotePeerRepo.On("FindOne", peerId.String()).Once().Return(remotePeerRow, nil)

		efs.OnNodeAdded(peerId)
	})
}
