package whatsapp

import (
	"log"
	"strings"
	"wppserver/pkg/config"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	
	uuid "github.com/google/uuid"
	_ "github.com/mdp/qrterminal/v3"
	_ "github.com/skip2/go-qrcode"
	"google.golang.org/protobuf/proto"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/appstate"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type Device struct {
	Id             uuid.UUID `sql:"id;default:uuid_generate_v4()"`
	UserId         uuid.UUID `sql:"userid;default:uuid_generate_v4()"`
	Jid            types.JID `json:"jid"`
	Connected      string    `json:"connected"`
	QrCode         dataQrCode
	Client         *whatsmeow.Client
	EventHandlerID uint32
	FnEvent        func(d *Device, rawEvt interface{})
	FnUpdateDevice func(d *Device)
}

type dataQrCode struct {
	Base64QrCode string `json:"base64qrcode"`
	Expiration   string `json:"expiration"`
}

func (d *Device) ParseJID(arg string) (types.JID, bool) {
	if arg == "" {
		return types.NewJID("", types.DefaultUserServer), false
	}
	if arg[0] == '+' {
		arg = arg[1:]
	}

	// Basic only digit check for recipient phone number, we want to remove @server and .session
	phonenumber := ""
	phonenumber = strings.Split(arg, "@")[0]
	phonenumber = strings.Split(phonenumber, ".")[0]
	b := true
	for _, c := range phonenumber {
		if c < '0' || c > '9' {
			b = false
			break
		}
	}
	if b == false {
		log.Printf("Bad jid format, return empty: %v\n", phonenumber)
		recipient, _ := types.ParseJID("")
		return recipient, false
	}

	if !strings.ContainsRune(arg, '@') {
		return types.NewJID(arg, types.DefaultUserServer), true
	} else {
		recipient, err := types.ParseJID(arg)
		if err != nil {
			log.Fatalf("Invalid jid: %v\n", err)
			return recipient, false
		} else if recipient.User == "" {
			log.Fatalf("Invalid jid. No server specified: %v\n", err)
			return recipient, false
		}
		return recipient, true
	}
}

func (d *Device) StartClient(container *sqlstore.Container) *whatsmeow.Client {
	var deviceStore *store.Device
	var err error

	log.Printf("Starting connection to Whatsapp:  %v\n", d.UserId)

	if d.Jid.String() != "" {
		deviceStore, err = container.GetDevice(d.Jid)
		if err != nil {
			log.Panic(err)
		}
	}

	if deviceStore == nil {
		log.Print("No valid jid found. Creating new device")
		deviceStore = container.NewDevice()
	}

	store.DeviceProps.PlatformType = waProto.DeviceProps_UNKNOWN.Enum()
	store.DeviceProps.Os = proto.String("Mac OS 10")

	clientLog := waLog.Stdout("Client", "DEBUG", true)

	d.Client = whatsmeow.NewClient(deviceStore, clientLog)
	d.EventHandlerID = d.Client.AddEventHandler(d.eventHandlerDevices)
	return d.Client
}

func (d *Device) eventHandlerDevices(rawEvt interface{}) {
	switch evt := rawEvt.(type) {
	case *events.AppStateSyncComplete:
		if len(d.Client.Store.PushName) > 0 && evt.Name == appstate.WAPatchCriticalBlock {
			err := d.Client.SendPresence(types.PresenceAvailable)
			if err != nil {
				log.Printf("Failed to send available presence: %v\n", err)
			} else {
				log.Print("Marked self as available")
			}
		}
	case *events.Connected, *events.PushNameSetting:
		if len(d.Client.Store.PushName) == 0 {
			return
		}
		err := d.Client.SendPresence(types.PresenceAvailable)
		if err != nil {
			log.Printf("Failed to send available presence: %v\n", err)
		} else {
			log.Print("Marked self as available")
		}
	case *events.PairSuccess:
		d.Jid = evt.ID
		log.Printf("New jid: %v\n", d.Jid)
		d.FnUpdateDevice(d)
	}

	//d.FnEvent(d, rawEvt)
}

func (d *Device) DisconnectClient() {
	log.Printf("Disconnecting to Whatsapp:  %v\n", d)

	if d.Client.IsConnected() {
		d.Client.Disconnect()
	}
}

type Devices struct {
	AllDevices map[uuid.UUID]*Device
	Container  *sqlstore.Container
}

func NewDevices() *Devices {
	config := config.GetConfig()
	myDevices := Devices{}
	myDevices.AllDevices = make(map[uuid.UUID]*Device)

	dbLog := waLog.Stdout("Database", "DEBUG", false)
	myDevices.Container, _ = sqlstore.New(config.DB.Dialect, config.DB.DbURI, dbLog)

	return &myDevices
}

func (ds *Devices) RegisterNew(dUUID uuid.UUID) (*Device, bool) {
	var err error

	myDevice := Device{}

	myDevice.Id, err = uuid.NewRandom()
	myDevice.UserId = dUUID

	if err != nil {
		panic(err)
	}

	ds.Add(&myDevice)

	return ds.Get(dUUID)
}

func (ds *Devices) Get(did uuid.UUID) (*Device, bool) {
	if _, ok := ds.AllDevices[did]; ok {
		return ds.AllDevices[did], true
	}
	return nil, false
}

func (ds *Devices) GetBy(sType string, value string) (*Device, bool) {
	for _, d := range ds.AllDevices {
		if sType == "phone" {
			phone := strings.Split(d.Jid.String(), "@")[0]
			phone = strings.Split(phone, ".")[0]

			if value == phone {
				return d, true
			}
		} else if sType == "jid" {
			if value == d.Jid.String() {
				return d, true
			}
		}
	}
	return nil, false
}

func (ds *Devices) Add(d *Device) {
	ds.AllDevices[d.UserId] = d
}

func (ds *Devices) Remove(did uuid.UUID) bool {
	d := ds.AllDevices[did]
	d.DisconnectClient()
	d = nil
	delete(ds.AllDevices, did)

	if _, ok := ds.AllDevices[did]; ok {
		return false
	}

	return true
}
