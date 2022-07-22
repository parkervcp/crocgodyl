package crocgodyl

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

type Limits struct {
	Memory      int64  `json:"memory"`
	Swap        int64  `json:"swap"`
	Disk        int64  `json:"disk"`
	IO          int64  `json:"io"`
	CPU         int64  `json:"cpu"`
	Threads     string `json:"threads"`
	OOMDisabled bool   `json:"oom_disabled"`
}

type FeatureLimits struct {
	Allocations int `json:"allocations"`
	Backups     int `json:"backups"`
	Databases   int `json:"databases"`
}

type ClientServer struct {
	ServerOwner bool   `json:"server_owner"`
	Identifier  string `json:"identifier"`
	UUID        string `json:"uuid"`
	InternalID  int    `json:"internal_id"`
	Name        string `json:"name"`
	Node        string `json:"node"`
	SFTP        struct {
		IP   string `json:"ip"`
		Port int64  `json:"port"`
	} `json:"sftp_details"`
	Description   string        `json:"description"`
	Limits        Limits        `json:"limits"`
	Invocation    string        `json:"invocation"`
	DockerImage   string        `json:"docker_image"`
	EggFeatures   []string      `json:"egg_features"`
	FeatureLimits FeatureLimits `json:"feature_limits"`
	Status        string        `json:"status"`
	Suspended     bool          `json:"is_suspended"`
	Installing    bool          `json:"is_installing"`
	Transferring  bool          `json:"is_transferring"`
}

func (c *Client) GetServers() ([]*ClientServer, error) {
	req := c.newRequest("GET", "", nil)
	res, err := c.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model struct {
		Data []struct {
			Attributes *ClientServer `json:"attributes"`
		} `json:"data"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	servers := make([]*ClientServer, 0, len(model.Data))
	for _, s := range model.Data {
		servers = append(servers, s.Attributes)
	}

	return servers, nil
}

func (c *Client) GetServer(identifier string) (*ClientServer, error) {
	req := c.newRequest("GET", "/servers/"+identifier, nil)
	res, err := c.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model struct {
		Attributes ClientServer `json:"attributes"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	return &model.Attributes, err
}

type WebSocketAuth struct {
	Socket string `json:"socket"`
	Token  string `json:"token"`
}

func (c *Client) GetServerWebSocket(identifier string) (*WebSocketAuth, error) {
	req := c.newRequest("GET", fmt.Sprintf("/servers/%s/websocket", identifier), nil)
	res, err := c.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model struct {
		Data WebSocketAuth `json:"data"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	return &model.Data, nil
}

type Resources struct {
	MemoryBytes    int64 `json:"memory_bytes"`
	DiskBytes      int64 `json:"disk_bytes"`
	CPUAbsolute    int64 `json:"cpu_absolute"`
	NetworkRxBytes int64 `json:"network_rx_bytes"`
	NetworkTxBytes int64 `json:"network_tx_bytes"`
	Uptime         int64 `json:"uptime"`
}

type Stats struct {
	State     string    `json:"current_state,omitempty"`
	Suspended bool      `json:"is_suspended"`
	Resources Resources `json:"resources"`
}

func (c *Client) ServerStatistics(identifier string) (*Stats, error) {
	req := c.newRequest("GET", fmt.Sprintf("/servers/%s/resources", identifier), nil)
	res, err := c.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model struct {
		Attributes Stats `json:"attributes"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	return &model.Attributes, nil
}

func (c *Client) SendServerCommand(identifier, command string) error {
	data, _ := json.Marshal(map[string]string{"command": command})
	body := bytes.Buffer{}
	body.Write(data)

	req := c.newRequest("POST", fmt.Sprintf("/servers/%s/command", identifier), &body)
	res, err := c.Http.Do(req)
	if err != nil {
		return err
	}

	_, err = validate(res)
	return err
}

func (c *Client) SetServerPowerState(identifier, state string) error {
	data, _ := json.Marshal(map[string]string{"signal": state})
	body := bytes.Buffer{}
	body.Write(data)

	req := c.newRequest("POST", fmt.Sprintf("/servers/%s/power", identifier), &body)
	res, err := c.Http.Do(req)
	if err != nil {
		return err
	}

	_, err = validate(res)
	return err
}

type ClientDatabase struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Host     struct {
		Address string `json:"address"`
		Port    int64  `json:"port"`
	} `json:"host"`
	ConnectionsFrom string `json:"connections_from"`
	MaxConnections  int    `json:"max_connections"`
}

func (c *Client) GetServerDatabases(identifier string) ([]*ClientDatabase, error) {
	req := c.newRequest("POST", fmt.Sprintf("/servers/%s/command", identifier), nil)
	res, err := c.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model struct {
		Data []struct {
			Attributes *ClientDatabase `json:"attributes"`
		} `json:"data"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	dbs := make([]*ClientDatabase, 0, len(model.Data))
	for _, d := range model.Data {
		dbs = append(dbs, d.Attributes)
	}

	return dbs, nil
}

func (c *Client) CreateDatabase(identifier, remote, database string) (*ClientDatabase, error) {
	data, _ := json.Marshal(map[string]string{"remote": remote, "database": database})
	body := bytes.Buffer{}
	body.Write(data)

	req := c.newRequest("POST", fmt.Sprintf("/servers/%s/databases", identifier), &body)
	res, err := c.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model struct {
		Attributes ClientDatabase `json:"attributes"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	return &model.Attributes, nil
}

func (c *Client) RotateDatabasePassword(identifier, id string) (*ClientDatabase, error) {
	req := c.newRequest("POST", fmt.Sprintf("/servers/%s/databases/%s/rotate-password", identifier, id), nil)
	res, err := c.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model struct {
		Attributes ClientDatabase `json:"attributes"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	return &model.Attributes, nil
}

func (c *Client) DeleteDatabase(identifier, id string) error {
	req := c.newRequest("DELETE", fmt.Sprintf("/servers/%s/databases/%s", identifier, id), nil)
	res, err := c.Http.Do(req)
	if err != nil {
		return err
	}

	_, err = validate(res)
	return err
}

type File struct {
	Name       string     `json:"name"`
	Mode       string     `json:"mode"`
	ModeBits   string     `json:"mode_bits"`
	Size       int64      `json:"size"`
	IsFile     bool       `json:"is_file"`
	IsSymlink  bool       `json:"is_symlink"`
	MimeType   string     `json:"mimetype"`
	CreatedAt  *time.Time `json:"created_at"`
	ModifiedAt *time.Time `json:"modified_at,omitempty"`
}

func (c *Client) GetServerFiles(identififer, root string) ([]*File, error) {
	req := c.newRequest("GET", fmt.Sprintf("/servers/%s/files/list?directory=%s", identififer, url.PathEscape(root)), nil)
	res, err := c.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model struct {
		Data []struct {
			Attributes *File `json:"attributes"`
		} `json:"data"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	files := make([]*File, 0, len(model.Data))
	for _, f := range model.Data {
		files = append(files, f.Attributes)
	}

	return files, nil
}

func (c *Client) GetServerFileContents(identifier, file string) ([]byte, error) {
	req := c.newRequest("GET", fmt.Sprintf("/servers/%s/files/contents?file=%s", identifier, url.PathEscape(file)), nil)
	req.Header.Set("Accept", "application/json,text/plain")

	res, err := c.Http.Do(req)
	if err != nil {
		return nil, err
	}

	return validate(res)
}

type Downloader struct {
	client *Client
	Name   string
	Path   string
	url    string
}

func (d *Downloader) Client() *Client {
	return d.client
}

func (d *Downloader) URL() string {
	return d.url
}

func (d *Downloader) Execute() error {
	info, err := os.Stat(d.Path)
	if err == nil {
		if !info.IsDir() {
			return errors.New("refusing to overwrite existing file path")
		}
	}

	res, err := d.client.Http.Get(d.URL())
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("recieved an unexpected response: %s", res.Status)
	}

	file, err := os.OpenFile(d.Name, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	io.Copy(file, res.Body)
	return nil
}

func (c *Client) DownloadServerFile(identifier, file string) (*Downloader, error) {
	files, err := c.GetServerFiles(identifier, "/")
	if err != nil {
		return nil, err
	}

	_, name := filepath.Split(file)
	for _, f := range files {
		if f.Name == name {
			if f.MimeType == "inode/directory" {
				return nil, errors.New("cannot download a directory")
			}

			break
		}
	}

	req := c.newRequest("GET", fmt.Sprintf("/servers/%s/files/download?file=%s", identifier, url.PathEscape(file)), nil)
	res, err := c.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model struct {
		Attributes struct {
			URL string `json:"url"`
		} `json:"attributes"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	path, _ := url.PathUnescape(file)
	dl := &Downloader{
		client: c,
		Name:   name,
		Path:   path,
		url:    model.Attributes.URL,
	}

	return dl, nil
}

type RenameDescriptor struct {
	Root  string `json:"root"`
	Files []struct {
		From string `json:"from"`
		To   string `json:"to"`
	} `json:"files"`
}

func (c *Client) RenameServerFiles(identifier string, files RenameDescriptor) error {
	data, _ := json.Marshal(files)
	body := bytes.Buffer{}
	body.Write(data)

	req := c.newRequest("PUT", fmt.Sprintf("/servers/%s/files/rename", identifier), &body)
	res, err := c.Http.Do(req)
	if err != nil {
		return err
	}

	_, err = validate(res)
	return err
}

func (c *Client) CopyServerFile(identifier, location string) error {
	data, _ := json.Marshal(map[string]string{"location": location})
	body := bytes.Buffer{}
	body.Write(data)

	req := c.newRequest("POST", fmt.Sprintf("/servers/%s/files/copy", identifier), &body)
	res, err := c.Http.Do(req)
	if err != nil {
		return err
	}

	_, err = validate(res)
	return err
}

func (c *Client) WriteServerFileBytes(identifier, name, header string, content []byte) error {
	body := bytes.Buffer{}
	body.Write(content)

	req := c.newRequest("POST", fmt.Sprintf("/servers/%s/files/write?file=%s", identifier, url.PathEscape(name)), &body)
	req.Header.Set("Content-Type", header)
	res, err := c.Http.Do(req)
	if err != nil {
		return err
	}

	_, err = validate(res)
	return err
}

func (c *Client) WriteServerFile(identifier, name, content string) error {
	return c.WriteServerFileBytes(identifier, name, "text/plain", []byte(content))
}

type CompressDescriptor struct {
	Root  string   `json:"root"`
	Files []string `json:"files"`
}

func (c *Client) CompressServerFiles(identifier string, files CompressDescriptor) error {
	data, _ := json.Marshal(files)
	body := bytes.Buffer{}
	body.Write(data)

	req := c.newRequest("POST", fmt.Sprintf("/servers/%s/files/compress", identifier), &body)
	res, err := c.Http.Do(req)
	if err != nil {
		return err
	}

	_, err = validate(res)
	return err
}

type DecompressDescriptor struct {
	Root string `json:"root"`
	File string `json:"file"`
}

func (c *Client) DecompressServerFile(identifier string, file DecompressDescriptor) error {
	data, _ := json.Marshal(file)
	body := bytes.Buffer{}
	body.Write(data)

	req := c.newRequest("POST", fmt.Sprintf("/servers/%s/files/decompress", identifier), &body)
	res, err := c.Http.Do(req)
	if err != nil {
		return err
	}

	_, err = validate(res)
	return err
}

type DeleteFilesDescriptor struct {
	Root  string   `json:"root"`
	Files []string `json:"files"`
}

func (c *Client) DeleteServerFiles(identifier string, files DeleteFilesDescriptor) error {
	data, _ := json.Marshal(files)
	body := bytes.Buffer{}
	body.Write(data)

	req := c.newRequest("POST", fmt.Sprintf("/servers/%s/files/delete", identifier), &body)
	res, err := c.Http.Do(req)
	if err != nil {
		return err
	}

	_, err = validate(res)
	return err
}

type CreateFolderDescriptor struct {
	Root string `json:"root"`
	Name string `json:"name"`
}

func (c *Client) CreateServerFileFolder(identifier string, file CreateFolderDescriptor) error {
	data, _ := json.Marshal(file)
	body := bytes.Buffer{}
	body.Write(data)

	req := c.newRequest("POST", fmt.Sprintf("/servers/%s/files/create-folder", identifier), &body)
	res, err := c.Http.Do(req)
	if err != nil {
		return err
	}

	_, err = validate(res)
	return err
}

type ChmodDescriptor struct {
	Root  string `json:"root"`
	Files []struct {
		File string `json:"file"`
		Mode uint32 `json:"mode"`
	} `json:"files"`
}

func (c *Client) ChmodServerFiles(identifier string, files ChmodDescriptor) error {
	data, _ := json.Marshal(files)
	body := bytes.Buffer{}
	body.Write(data)

	req := c.newRequest("POST", fmt.Sprintf("/servers/%s/files/chmod", identifier), &body)
	res, err := c.Http.Do(req)
	if err != nil {
		return err
	}

	_, err = validate(res)
	return err
}

type PullDescriptor struct {
	URL        string `json:"url"`
	Directory  string `json:"directory,omitempty"`
	Filename   string `json:"filename,omitempty"`
	UseHeader  bool   `json:"use_header,omitempty"`
	Foreground bool   `json:"foreground,omitempty"`
}

func (c *Client) PullServerFile(identifier string, file PullDescriptor) error {
	data, _ := json.Marshal(file)
	body := bytes.Buffer{}
	body.Write(data)

	req := c.newRequest("POST", fmt.Sprintf("/servers/%s/files/pull", identifier), &body)
	res, err := c.Http.Do(req)
	if err != nil {
		return err
	}

	_, err = validate(res)
	return err
}

type Uploader struct {
	client *Client
	url    string
	Path   string
}

func (u *Uploader) Client() *Client {
	return u.client
}

func (u *Uploader) URL() string {
	return u.url
}

func (u *Uploader) Execute() error {
	if u.Path == "" {
		return errors.New("no file path has been specified")
	}

	info, err := os.Stat(u.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("file path does not exist")
		}

		return err
	}

	if info.IsDir() {
		return errors.New("path must go to a file not a directory")
	}

	file, err := os.Open(u.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	body := bytes.Buffer{}
	writer := multipart.NewWriter(&body)
	part, _ := writer.CreateFormFile("files", info.Name())
	io.Copy(part, file)
	writer.Close()

	res, err := u.client.Http.Post(u.URL(), writer.FormDataContentType(), &body)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("recieved an unexpected response: %s", res.Status)
	}

	return nil
}

func (c *Client) UploadServerFile(identifier, path string) (*Uploader, error) {
	req := c.newRequest("GET", fmt.Sprintf("/servers/%s/files/upload", identifier), nil)
	res, err := c.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model struct {
		Attributes struct {
			URL string `json:"url"`
		} `json:"attributes"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	up := &Uploader{client: c, url: model.Attributes.URL}
	return up, nil
}
