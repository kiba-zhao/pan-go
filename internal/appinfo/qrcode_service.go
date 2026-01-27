package appinfo

import (
	"bytes"
	"errors"
	"image/color"
	"io"
	"net/url"
	"pan/internal/peer"
	"pan/internal/settings"
	"sync"
	"time"

	"github.com/skip2/go-qrcode"
)

var ErrQRCodeUnavailable = errors.New("appinfo.QRCodeService Error: Unavailable")

type QRCodeReader struct {
	io.ReadSeeker
	name    string
	modTime time.Time
}

func (reader *QRCodeReader) Name() string {
	return reader.name
}

func (reader *QRCodeReader) ModTime() time.Time {
	return reader.modTime
}

type QRCodeService struct {
	SettingsConfigurer settings.SettingsConfigurer
	SecurityConfigurer settings.SecurityConfigurer
	SettingsService    settings.SettingsExternalService

	peerId  peer.PeerID
	name    string
	modTime time.Time
	locker  sync.Mutex
}

func (service *QRCodeService) Load(size int) (*QRCodeReader, error) {
	if service.SecurityConfigurer == nil || service.SettingsService == nil || service.SettingsConfigurer == nil {
		return nil, ErrQRCodeUnavailable
	}

	settingsConf := service.SettingsConfigurer.Config()
	if settingsConf == nil {
		return nil, ErrQRCodeUnavailable
	}

	securityConf := service.SecurityConfigurer.Config()
	if securityConf == nil {
		return nil, ErrQRCodeUnavailable
	}

	settings, err := service.SettingsService.Load()
	if err != nil {
		return nil, err
	}

	peerID := securityConf.PeerID()
	name := settings.Name
	var modTime time.Time
	service.locker.Lock()
	if !bytes.Equal(service.peerId, peerID) || service.name != name {
		service.peerId = peerID
		service.name = name
		modTime = time.Now()
	} else {
		modTime = service.modTime
	}
	service.locker.Unlock()

	var qrcode_ *qrcode.QRCode
	var contentColor color.Color
	if len(peerID) > 0 && len(name) > 0 {
		qrcode_, err = qrcode.New(generateQRCodeValue(peerID, name), qrcode.Medium)
		contentColor = color.Black
	} else {
		qrcode_, err = qrcode.New("", qrcode.Medium)
		contentColor = color.RGBA{
			R: 255,
			G: 0,
			B: 0,
			A: 255,
		}
	}

	if err != nil {
		return nil, err
	}

	qrcode_.BackgroundColor = color.White
	qrcode_.ForegroundColor = contentColor
	qrcodeSize := size
	if qrcodeSize == 0 {
		qrcodeSize = 256
	}
	qrcodeBytes, err := qrcode_.PNG(qrcodeSize)
	if err != nil {
		return nil, err
	}

	reader := &QRCodeReader{}
	reader.ReadSeeker = bytes.NewReader(qrcodeBytes)
	reader.modTime = modTime

	if len(name) > 0 {
		reader.name = name + ".png"
	} else {
		reader.name = settingsConf.HostName() + ".png"
	}

	return reader, nil
}

func generateQRCodeValue(peerID []byte, name string) string {
	v := url.Values{}
	v.Set("peerId", peer.EncodePeerID(peerID))
	v.Set("name", name)

	return "pan://appinfo?" + v.Encode()
}
