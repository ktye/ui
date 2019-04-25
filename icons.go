package ui

/* TODO port to v2

import (
	"image"
	"image/draw"

	"golang.org/x/exp/shiny/iconvg"
)

func init() {
	iconVgs = make(map[string][]byte)
	icons = make(map[icon]image.Image)
}

// RegisterIcon stores an icon in iconvg format for a given name.
// See github.com/golang/exp/blob/master/shiny/materialdesign/icons/data_test.go
// for the full list.
//
// Example: RegisterIcon("ActionAccoutCircle", icons.ActionAccountCircle)
func RegisterIcon(name string, data []byte) {
	iconVgs[name] = data
}

// IconVgs stores icons in iconvg format under a given name.
var iconVgs map[string][]byte

// Icon is used as a key to the map of pre-rendered images for a given size.
type icon struct {
	name string
	size int
}

// Icons stores pre-rendered images for the given icon name an size.
var icons map[icon]image.Image

// GetIcon returns an icon image from cache or renders the icon on demand.
// If it does not exist, a grey image is returned.
func getIcon(w *Window, id icon) image.Image {
	if id.name == "" || id.size <= 0 {
		return unknownIcon(w.FontHeight())
	}
	ic, ok := icons[id]
	if ok {
		return ic
	}

	data, ok := iconVgs[id.name]
	if ok == false {
		return unknownIcon(id.size)
	}

	img := image.NewAlpha(image.Rect(0, 0, id.size, id.size))
	z := &iconvg.Rasterizer{}
	z.SetDstImage(img, img.Bounds(), draw.Src)
	if err := iconvg.Decode(z, data, nil); err != nil {
		return unknownIcon(id.size)
	}
	icons[id] = img
	return img
}

func unknownIcon(size int) image.Image {
	return image.NewGray(rect(image.Point{size, size}))
}

*/

/* List of icons in golang.org/x/exp/shiny/materialdesign/icons
Action3DRotation
ActionAccessibility
ActionAccessible
ActionAccountBalance
ActionAccountBalanceWallet
ActionAccountBox
ActionAccountCircle
ActionAddShoppingCart
ActionAlarm
ActionAlarmAdd
ActionAlarmOff
ActionAlarmOn
ActionAllOut
ActionAndroid
ActionAnnouncement
ActionAspectRatio
ActionAssessment
ActionAssignment
ActionAssignmentInd
ActionAssignmentLate
ActionAssignmentReturn
ActionAssignmentReturned
ActionAssignmentTurnedIn
ActionAutorenew
ActionBackup
ActionBook
ActionBookmark
ActionBookmarkBorder
ActionBugReport
ActionBuild
ActionCached
ActionCameraEnhance
ActionCardGiftcard
ActionCardMembership
ActionCardTravel
ActionChangeHistory
ActionCheckCircle
ActionChromeReaderMode
ActionClass
ActionCode
ActionCompareArrows
ActionCopyright
ActionCreditCard
ActionDashboard
ActionDateRange
ActionDelete
ActionDeleteForever
ActionDescription
ActionDNS
ActionDone
ActionDoneAll
ActionDonutLarge
ActionDonutSmall
ActionEject
ActionEuroSymbol
ActionEvent
ActionEventSeat
ActionExitToApp
ActionExplore
ActionExtension
ActionFace
ActionFavorite
ActionFavoriteBorder
ActionFeedback
ActionFindInPage
ActionFindReplace
ActionFingerprint
ActionFlightLand
ActionFlightTakeoff
ActionFlipToBack
ActionFlipToFront
ActionGTranslate
ActionGavel
ActionGetApp
ActionGIF
ActionGrade
ActionGroupWork
ActionHelp
ActionHelpOutline
ActionHighlightOff
ActionHistory
ActionHome
ActionHourglassEmpty
ActionHourglassFull
ActionHTTP
ActionHTTPS
ActionImportantDevices
ActionInfo
ActionInfoOutline
ActionInput
ActionInvertColors
ActionLabel
ActionLabelOutline
ActionLanguage
ActionLaunch
ActionLightbulbOutline
ActionLineStyle
ActionLineWeight
ActionList
ActionLock
ActionLockOpen
ActionLockOutline
ActionLoyalty
ActionMarkUnreadMailbox
ActionMotorcycle
ActionNoteAdd
ActionOfflinePin
ActionOpacity
ActionOpenInBrowser
ActionOpenInNew
ActionOpenWith
ActionPageview
ActionPanTool
ActionPayment
ActionPermCameraMic
ActionPermContactCalendar
ActionPermDataSetting
ActionPermDeviceInformation
ActionPermIdentity
ActionPermMedia
ActionPermPhoneMsg
ActionPermScanWiFi
ActionPets
ActionPictureInPicture
ActionPictureInPictureAlt
ActionPlayForWork
ActionPolymer
ActionPowerSettingsNew
ActionPregnantWoman
ActionPrint
ActionQueryBuilder
ActionQuestionAnswer
ActionReceipt
ActionRecordVoiceOver
ActionRedeem
ActionRemoveShoppingCart
ActionReorder
ActionReportProblem
ActionRestore
ActionRestorePage
ActionRoom
ActionRoundedCorner
ActionRowing
ActionSchedule
ActionSearch
ActionSettings
ActionSettingsApplications
ActionSettingsBackupRestore
ActionSettingsBluetooth
ActionSettingsBrightness
ActionSettingsCell
ActionSettingsEthernet
ActionSettingsInputAntenna
ActionSettingsInputComponent
ActionSettingsInputComposite
ActionSettingsInputHDMI
ActionSettingsInputSVideo
ActionSettingsOverscan
ActionSettingsPhone
ActionSettingsPower
ActionSettingsRemote
ActionSettingsVoice
ActionShop
ActionShopTwo
ActionShoppingBasket
ActionShoppingCart
ActionSpeakerNotes
ActionSpeakerNotesOff
ActionSpellcheck
ActionStarRate
ActionStars
ActionStore
ActionSubject
ActionSupervisorAccount
ActionSwapHoriz
ActionSwapVert
ActionSwapVerticalCircle
ActionSystemUpdateAlt
ActionTab
ActionTabUnselected
ActionTheaters
ActionThumbDown
ActionThumbUp
ActionThumbsUpDown
ActionTimeline
ActionTOC
ActionToday
ActionToll
ActionTouchApp
ActionTrackChanges
ActionTranslate
ActionTrendingDown
ActionTrendingFlat
ActionTrendingUp
ActionTurnedIn
ActionTurnedInNot
ActionUpdate
ActionVerifiedUser
ActionViewAgenda
ActionViewArray
ActionViewCarousel
ActionViewColumn
ActionViewDay
ActionViewHeadline
ActionViewList
ActionViewModule
ActionViewQuilt
ActionViewStream
ActionViewWeek
ActionVisibility
ActionVisibilityOff
ActionWatchLater
ActionWork
ActionYoutubeSearchedFor
ActionZoomIn
ActionZoomOut
AlertAddAlert
AlertError
AlertErrorOutline
AlertWarning
AVAddToQueue
AVAirplay
AVAlbum
AVArtTrack
AVAVTimer
AVBrandingWatermark
AVCallToAction
AVClosedCaption
AVEqualizer
AVExplicit
AVFastForward
AVFastRewind
AVFeaturedPlayList
AVFeaturedVideo
AVFiberDVR
AVFiberManualRecord
AVFiberNew
AVFiberPin
AVFiberSmartRecord
AVForward10
AVForward30
AVForward5
AVGames
AVHD
AVHearing
AVHighQuality
AVLibraryAdd
AVLibraryBooks
AVLibraryMusic
AVLoop
AVMic
AVMicNone
AVMicOff
AVMovie
AVMusicVideo
AVNewReleases
AVNotInterested
AVNote
AVPause
AVPauseCircleFilled
AVPauseCircleOutline
AVPlayArrow
AVPlayCircleFilled
AVPlayCircleOutline
AVPlaylistAdd
AVPlaylistAddCheck
AVPlaylistPlay
AVQueue
AVQueueMusic
AVQueuePlayNext
AVRadio
AVRecentActors
AVRemoveFromQueue
AVRepeat
AVRepeatOne
AVReplay
AVReplay10
AVReplay30
AVReplay5
AVShuffle
AVSkipNext
AVSkipPrevious
AVSlowMotionVideo
AVSnooze
AVSortByAlpha
AVStop
AVSubscriptions
AVSubtitles
AVSurroundSound
AVVideoCall
AVVideoLabel
AVVideoLibrary
AVVideocam
AVVideocamOff
AVVolumeDown
AVVolumeMute
AVVolumeOff
AVVolumeUp
AVWeb
AVWebAsset
CommunicationBusiness
CommunicationCall
CommunicationCallEnd
CommunicationCallMade
CommunicationCallMerge
CommunicationCallMissed
CommunicationCallMissedOutgoing
CommunicationCallReceived
CommunicationCallSplit
CommunicationChat
CommunicationChatBubble
CommunicationChatBubbleOutline
CommunicationClearAll
CommunicationComment
CommunicationContactMail
CommunicationContactPhone
CommunicationContacts
CommunicationDialerSIP
CommunicationDialpad
CommunicationEmail
CommunicationForum
CommunicationImportContacts
CommunicationImportExport
CommunicationInvertColorsOff
CommunicationLiveHelp
CommunicationLocationOff
CommunicationLocationOn
CommunicationMailOutline
CommunicationMessage
CommunicationNoSIM
CommunicationPhone
CommunicationPhoneLinkErase
CommunicationPhoneLinkLock
CommunicationPhoneLinkRing
CommunicationPhoneLinkSetup
CommunicationPortableWiFiOff
CommunicationPresentToAll
CommunicationRingVolume
CommunicationRSSFeed
CommunicationScreenShare
CommunicationSpeakerPhone
CommunicationStayCurrentLandscape
CommunicationStayCurrentPortrait
CommunicationStayPrimaryLandscape
CommunicationStayPrimaryPortrait
CommunicationStopScreenShare
CommunicationSwapCalls
CommunicationTextSMS
CommunicationVoicemail
CommunicationVPNKey
ContentAdd
ContentAddBox
ContentAddCircle
ContentAddCircleOutline
ContentArchive
ContentBackspace
ContentBlock
ContentClear
ContentContentCopy
ContentContentCut
ContentContentPaste
ContentCreate
ContentDeleteSweep
ContentDrafts
ContentFilterList
ContentFlag
ContentFontDownload
ContentForward
ContentGesture
ContentInbox
ContentLink
ContentLowPriority
ContentMail
ContentMarkUnread
ContentMoveToInbox
ContentNextWeek
ContentRedo
ContentRemove
ContentRemoveCircle
ContentRemoveCircleOutline
ContentReply
ContentReplyAll
ContentReport
ContentSave
ContentSelectAll
ContentSend
ContentSort
ContentTextFormat
ContentUnarchive
ContentUndo
ContentWeekend
DeviceAccessAlarm
DeviceAccessAlarms
DeviceAccessTime
DeviceAddAlarm
DeviceAirplaneModeActive
DeviceAirplaneModeInactive
DeviceBattery20
DeviceBattery30
DeviceBattery50
DeviceBattery60
DeviceBattery80
DeviceBattery90
DeviceBatteryAlert
DeviceBatteryCharging20
DeviceBatteryCharging30
DeviceBatteryCharging50
DeviceBatteryCharging60
DeviceBatteryCharging80
DeviceBatteryCharging90
DeviceBatteryChargingFull
DeviceBatteryFull
DeviceBatteryStd
DeviceBatteryUnknown
DeviceBluetooth
DeviceBluetoothConnected
DeviceBluetoothDisabled
DeviceBluetoothSearching
DeviceBrightnessAuto
DeviceBrightnessHigh
DeviceBrightnessLow
DeviceBrightnessMedium
DeviceDataUsage
DeviceDeveloperMode
DeviceDevices
DeviceDVR
DeviceGPSFixed
DeviceGPSNotFixed
DeviceGPSOff
DeviceGraphicEq
DeviceLocationDisabled
DeviceLocationSearching
DeviceNetworkCell
DeviceNetworkWiFi
DeviceNFC
DeviceScreenLockLandscape
DeviceScreenLockPortrait
DeviceScreenLockRotation
DeviceScreenRotation
DeviceSDStorage
DeviceSettingsSystemDaydream
DeviceSignalCellular0Bar
DeviceSignalCellular1Bar
DeviceSignalCellular2Bar
DeviceSignalCellular3Bar
DeviceSignalCellular4Bar
DeviceSignalCellularConnectedNoInternet0Bar
DeviceSignalCellularConnectedNoInternet1Bar
DeviceSignalCellularConnectedNoInternet2Bar
DeviceSignalCellularConnectedNoInternet3Bar
DeviceSignalCellularConnectedNoInternet4Bar
DeviceSignalCellularNoSIM
DeviceSignalCellularNull
DeviceSignalCellularOff
DeviceSignalWiFi0Bar
DeviceSignalWiFi1Bar
DeviceSignalWiFi1BarLock
DeviceSignalWiFi2Bar
DeviceSignalWiFi2BarLock
DeviceSignalWiFi3Bar
DeviceSignalWiFi3BarLock
DeviceSignalWiFi4Bar
DeviceSignalWiFi4BarLock
DeviceSignalWiFiOff
DeviceStorage
DeviceUSB
DeviceWallpaper
DeviceWidgets
DeviceWiFiLock
DeviceWiFiTethering
EditorAttachFile
EditorAttachMoney
EditorBorderAll
EditorBorderBottom
EditorBorderClear
EditorBorderColor
EditorBorderHorizontal
EditorBorderInner
EditorBorderLeft
EditorBorderOuter
EditorBorderRight
EditorBorderStyle
EditorBorderTop
EditorBorderVertical
EditorBubbleChart
EditorDragHandle
EditorFormatAlignCenter
EditorFormatAlignJustify
EditorFormatAlignLeft
EditorFormatAlignRight
EditorFormatBold
EditorFormatClear
EditorFormatColorFill
EditorFormatColorReset
EditorFormatColorText
EditorFormatIndentDecrease
EditorFormatIndentIncrease
EditorFormatItalic
EditorFormatLineSpacing
EditorFormatListBulleted
EditorFormatListNumbered
EditorFormatPaint
EditorFormatQuote
EditorFormatShapes
EditorFormatSize
EditorFormatStrikethrough
EditorFormatTextDirectionLToR
EditorFormatTextDirectionRToL
EditorFormatUnderlined
EditorFunctions
EditorHighlight
EditorInsertChart
EditorInsertComment
EditorInsertDriveFile
EditorInsertEmoticon
EditorInsertInvitation
EditorInsertLink
EditorInsertPhoto
EditorLinearScale
EditorMergeType
EditorModeComment
EditorModeEdit
EditorMonetizationOn
EditorMoneyOff
EditorMultilineChart
EditorPieChart
EditorPieChartOutlined
EditorPublish
EditorShortText
EditorShowChart
EditorSpaceBar
EditorStrikethroughS
EditorTextFields
EditorTitle
EditorVerticalAlignBottom
EditorVerticalAlignCenter
EditorVerticalAlignTop
EditorWrapText
FileAttachment
FileCloud
FileCloudCircle
FileCloudDone
FileCloudDownload
FileCloudOff
FileCloudQueue
FileCloudUpload
FileCreateNewFolder
FileFileDownload
FileFileUpload
FileFolder
FileFolderOpen
FileFolderShared
HardwareCast
HardwareCastConnected
HardwareComputer
HardwareDesktopMac
HardwareDesktopWindows
HardwareDeveloperBoard
HardwareDeviceHub
HardwareDevicesOther
HardwareDock
HardwareGamepad
HardwareHeadset
HardwareHeadsetMic
HardwareKeyboard
HardwareKeyboardArrowDown
HardwareKeyboardArrowLeft
HardwareKeyboardArrowRight
HardwareKeyboardArrowUp
HardwareKeyboardBackspace
HardwareKeyboardCapslock
HardwareKeyboardHide
HardwareKeyboardReturn
HardwareKeyboardTab
HardwareKeyboardVoice
HardwareLaptop
HardwareLaptopChromebook
HardwareLaptopMac
HardwareLaptopWindows
HardwareMemory
HardwareMouse
HardwarePhoneAndroid
HardwarePhoneIPhone
HardwarePhoneLink
HardwarePhoneLinkOff
HardwarePowerInput
HardwareRouter
HardwareScanner
HardwareSecurity
HardwareSIMCard
HardwareSmartphone
HardwareSpeaker
HardwareSpeakerGroup
HardwareTablet
HardwareTabletAndroid
HardwareTabletMac
HardwareToys
HardwareTV
HardwareVideogameAsset
HardwareWatch
ImageAddAPhoto
ImageAddToPhotos
ImageAdjust
ImageAssistant
ImageAssistantPhoto
ImageAudiotrack
ImageBlurCircular
ImageBlurLinear
ImageBlurOff
ImageBlurOn
ImageBrightness1
ImageBrightness2
ImageBrightness3
ImageBrightness4
ImageBrightness5
ImageBrightness6
ImageBrightness7
ImageBrokenImage
ImageBrush
ImageBurstMode
ImageCamera
ImageCameraAlt
ImageCameraFront
ImageCameraRear
ImageCameraRoll
ImageCenterFocusStrong
ImageCenterFocusWeak
ImageCollections
ImageCollectionsBookmark
ImageColorLens
ImageColorize
ImageCompare
ImageControlPoint
ImageControlPointDuplicate
ImageCrop
ImageCrop169
ImageCrop32
ImageCrop54
ImageCrop75
ImageCropDIN
ImageCropFree
ImageCropLandscape
ImageCropOriginal
ImageCropPortrait
ImageCropRotate
ImageCropSquare
ImageDehaze
ImageDetails
ImageEdit
ImageExposure
ImageExposureNeg1
ImageExposureNeg2
ImageExposurePlus1
ImageExposurePlus2
ImageExposureZero
ImageFilter
ImageFilter1
ImageFilter2
ImageFilter3
ImageFilter4
ImageFilter5
ImageFilter6
ImageFilter7
ImageFilter8
ImageFilter9
ImageFilter9Plus
ImageFilterBAndW
ImageFilterCenterFocus
ImageFilterDrama
ImageFilterFrames
ImageFilterHDR
ImageFilterNone
ImageFilterTiltShift
ImageFilterVintage
ImageFlare
ImageFlashAuto
ImageFlashOff
ImageFlashOn
ImageFlip
ImageGradient
ImageGrain
ImageGridOff
ImageGridOn
ImageHDROff
ImageHDROn
ImageHDRStrong
ImageHDRWeak
ImageHealing
ImageImage
ImageImageAspectRatio
ImageISO
ImageLandscape
ImageLeakAdd
ImageLeakRemove
ImageLens
ImageLinkedCamera
ImageLooks
ImageLooks3
ImageLooks4
ImageLooks5
ImageLooks6
ImageLooksOne
ImageLooksTwo
ImageLoupe
ImageMonochromePhotos
ImageMovieCreation
ImageMovieFilter
ImageMusicNote
ImageNature
ImageNaturePeople
ImageNavigateBefore
ImageNavigateNext
ImagePalette
ImagePanorama
ImagePanoramaFishEye
ImagePanoramaHorizontal
ImagePanoramaVertical
ImagePanoramaWideAngle
ImagePhoto
ImagePhotoAlbum
ImagePhotoCamera
ImagePhotoFilter
ImagePhotoLibrary
ImagePhotoSizeSelectActual
ImagePhotoSizeSelectLarge
ImagePhotoSizeSelectSmall
ImagePictureAsPDF
ImagePortrait
ImageRemoveRedEye
ImageRotate90DegreesCCW
ImageRotateLeft
ImageRotateRight
ImageSlideshow
ImageStraighten
ImageStyle
ImageSwitchCamera
ImageSwitchVideo
ImageTagFaces
ImageTexture
ImageTimeLapse
ImageTimer
ImageTimer10
ImageTimer3
ImageTimerOff
ImageTonality
ImageTransform
ImageTune
ImageViewComfy
ImageViewCompact
ImageVignette
ImageWBAuto
ImageWBCloudy
ImageWBIncandescent
ImageWBIridescent
ImageWBSunny
MapsAddLocation
MapsBeenhere
MapsDirections
MapsDirectionsBike
MapsDirectionsBoat
MapsDirectionsBus
MapsDirectionsCar
MapsDirectionsRailway
MapsDirectionsRun
MapsDirectionsSubway
MapsDirectionsTransit
MapsDirectionsWalk
MapsEditLocation
MapsEVStation
MapsFlight
MapsHotel
MapsLayers
MapsLayersClear
MapsLocalActivity
MapsLocalAirport
MapsLocalATM
MapsLocalBar
MapsLocalCafe
MapsLocalCarWash
MapsLocalConvenienceStore
MapsLocalDining
MapsLocalDrink
MapsLocalFlorist
MapsLocalGasStation
MapsLocalGroceryStore
MapsLocalHospital
MapsLocalHotel
MapsLocalLaundryService
MapsLocalLibrary
MapsLocalMall
MapsLocalMovies
MapsLocalOffer
MapsLocalParking
MapsLocalPharmacy
MapsLocalPhone
MapsLocalPizza
MapsLocalPlay
MapsLocalPostOffice
MapsLocalPrintshop
MapsLocalSee
MapsLocalShipping
MapsLocalTaxi
MapsMap
MapsMyLocation
MapsNavigation
MapsNearMe
MapsPersonPin
MapsPersonPinCircle
MapsPinDrop
MapsPlace
MapsRateReview
MapsRestaurant
MapsRestaurantMenu
MapsSatellite
MapsStoreMallDirectory
MapsStreetView
MapsSubway
MapsTerrain
MapsTraffic
MapsTrain
MapsTram
MapsTransferWithinAStation
MapsZoomOutMap
NavigationApps
NavigationArrowBack
NavigationArrowDownward
NavigationArrowDropDown
NavigationArrowDropDownCircle
NavigationArrowDropUp
NavigationArrowForward
NavigationArrowUpward
NavigationCancel
NavigationCheck
NavigationChevronLeft
NavigationChevronRight
NavigationClose
NavigationExpandLess
NavigationExpandMore
NavigationFirstPage
NavigationFullscreen
NavigationFullscreenExit
NavigationLastPage
NavigationMenu
NavigationMoreHoriz
NavigationMoreVert
NavigationRefresh
NavigationSubdirectoryArrowLeft
NavigationSubdirectoryArrowRight
NavigationUnfoldLess
NavigationUnfoldMore
NotificationADB
NotificationAirlineSeatFlat
NotificationAirlineSeatFlatAngled
NotificationAirlineSeatIndividualSuite
NotificationAirlineSeatLegroomExtra
NotificationAirlineSeatLegroomNormal
NotificationAirlineSeatLegroomReduced
NotificationAirlineSeatReclineExtra
NotificationAirlineSeatReclineNormal
NotificationBluetoothAudio
NotificationConfirmationNumber
NotificationDiscFull
NotificationDoNotDisturb
NotificationDoNotDisturbAlt
NotificationDoNotDisturbOff
NotificationDoNotDisturbOn
NotificationDriveETA
NotificationEnhancedEncryption
NotificationEventAvailable
NotificationEventBusy
NotificationEventNote
NotificationFolderSpecial
NotificationLiveTV
NotificationMMS
NotificationMore
NotificationNetworkCheck
NotificationNetworkLocked
NotificationNoEncryption
NotificationOnDemandVideo
NotificationPersonalVideo
NotificationPhoneBluetoothSpeaker
NotificationPhoneForwarded
NotificationPhoneInTalk
NotificationPhoneLocked
NotificationPhoneMissed
NotificationPhonePaused
NotificationPower
NotificationPriorityHigh
NotificationRVHookup
NotificationSDCard
NotificationSIMCardAlert
NotificationSMS
NotificationSMSFailed
NotificationSync
NotificationSyncDisabled
NotificationSyncProblem
NotificationSystemUpdate
NotificationTapAndPlay
NotificationTimeToLeave
NotificationVibration
NotificationVoiceChat
NotificationVPNLock
NotificationWC
NotificationWiFi
PlacesACUnit
PlacesAirportShuttle
PlacesAllInclusive
PlacesBeachAccess
PlacesBusinessCenter
PlacesCasino
PlacesChildCare
PlacesChildFriendly
PlacesFitnessCenter
PlacesFreeBreakfast
PlacesGolfCourse
PlacesHotTub
PlacesKitchen
PlacesPool
PlacesRoomService
PlacesRVHookup
PlacesSmokeFree
PlacesSmokingRooms
PlacesSpa
SocialCake
SocialDomain
SocialGroup
SocialGroupAdd
SocialLocationCity
SocialMood
SocialMoodBad
SocialNotifications
SocialNotificationsActive
SocialNotificationsNone
SocialNotificationsOff
SocialNotificationsPaused
SocialPages
SocialPartyMode
SocialPeople
SocialPeopleOutline
SocialPerson
SocialPersonAdd
SocialPersonOutline
SocialPlusOne
SocialPoll
SocialPublic
SocialSchool
SocialSentimentDissatisfied
SocialSentimentNeutral
SocialSentimentSatisfied
SocialSentimentVeryDissatisfied
SocialSentimentVerySatisfied
SocialShare
SocialWhatsHot
ToggleCheckBox
ToggleCheckBoxOutlineBlank
ToggleIndeterminateCheckBox
ToggleRadioButtonChecked
ToggleRadioButtonUnchecked
ToggleStar
ToggleStarBorder
ToggleStarHalf
*/
