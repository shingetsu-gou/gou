/*
 * Copyright (c) 2015, Shinya Yagyu
 * All rights reserved.
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice,
 *    this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 * 3. Neither the name of the copyright holder nor the names of its
 *    contributors may be used to endorse or promote products derived from this
 *    software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
 * LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 */

package cfg

import (
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"sync"
	"time"

	"gopkg.in/ini.v1"
)

var NetworkMode string //port_opened,relay,upnp
var SaveRecord int64
var SaveSize int // It is not seconds, but number.
var GetRange int64
var SyncRange int64
var SaveRemoved int64
var DefaultPort int //DefaultPort is listening port
var MaxConnection int
var Docroot string
var LogDir string
var RunDir string
var FileDir string
var CacheDir string
var TemplateDir string
var SpamList string
var InitnodeList string
var FollowList string
var NodeAllowFile string
var NodeDenyFile string
var ReAdminStr string
var ReFriendStr string
var ReVisitorStr string
var ServerName string
var TagSize int
var RSSRange int64
var TopRecentRange int64
var RecentRange int64
var RecordLimit int
var ThreadPageSize int
var DefaultThumbnailSize string
var Enable2ch bool
var ForceThumbnail bool
var EnableProf bool
var HeavyMoon bool
var EnableEmbed bool
var RelayNumber int

// asis, md5, sha1, sha224, sha256, sha384, or sha512
//	cache_hash_method = "asis"
//others are not implemented for gou for now.

var Fmutex = &sync.RWMutex{}

//Version is one of Gou. it shoud be overwritten when building on travis.
var Version = "unstable"

//getIntValue gets int value from ini file.
func getIntValue(i *ini.File, section, key string, vdefault int) int {
	return i.Section(section).Key(key).MustInt(vdefault)
}

//getInt64Value gets int value from ini file.
func getInt64Value(i *ini.File, section, key string, vdefault int64) int64 {
	return i.Section(section).Key(key).MustInt64(vdefault)
}

//getStringValue gets string from ini file.
func getStringValue(i *ini.File, section, key string, vdefault string) string {
	return i.Section(section).Key(key).MustString(vdefault)
}

//getBoolValue gets bool value from ini file.
func getBoolValue(i *ini.File, section, key string, vdefault bool) bool {
	return i.Section(section).Key(key).MustBool(vdefault)
}

//getPathValue gets path from ini file.
func getRelativePathValue(i *ini.File, section, key, vdefault, docroot string) string {
	p := i.Section(section).Key(key).MustString(vdefault)
	h := p
	if !path.IsAbs(p) {
		h = path.Join(docroot, p)
	}
	return filepath.FromSlash(h)
}

//getPathValue gets path from ini file.
func getPathValue(i *ini.File, section, key string, vdefault string) string {
	p := i.Section(section).Key(key).MustString(vdefault)
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	h := p
	if !path.IsAbs(p) {
		h = path.Join(wd, p)
	}
	return filepath.FromSlash(h)
}

//init makes  config vars from the ini files and returns it.
func init() {
	files := []string{"file/saku.ini", "/usr/local/etc/saku/saku.ini", "/etc/saku/saku.ini"}
	usr, err := user.Current()
	if err == nil {
		files = append(files, usr.HomeDir+"/.saku/saku.ini")
	}
	i := ini.Empty()
	for _, f := range files {
		fs, err := os.Stat(f)
		if err != nil && !fs.IsDir() {
			log.Println("loading config from", f)
			if err := i.Append(f); err != nil {
				log.Fatal("cannot load ini files", f, "ignored")
			}
		}
	}
	initVariables(i)
}

//initVariables initializes some global and map vars.
func initVariables(i *ini.File) {
	DefaultPort = getIntValue(i, "Network", "port", 8000)
	NetworkMode = getStringValue(i, "Network", "mode", "port_opened") //port_opened,upnp,relay
	MaxConnection = getIntValue(i, "Network", "max_connection", 100)
	Docroot = getPathValue(i, "Path", "docroot", "./www")                                     //path from cwd
	RunDir = getRelativePathValue(i, "Path", "run_dir", "../run", Docroot)                    //path from docroot
	FileDir = getRelativePathValue(i, "Path", "file_dir", "../file", Docroot)                 //path from docroot
	CacheDir = getRelativePathValue(i, "Path", "cache_dir", "../cache", Docroot)              //path from docroot
	TemplateDir = getRelativePathValue(i, "Path", "template_dir", "../gou_template", Docroot) //path from docroot
	SpamList = getRelativePathValue(i, "Path", "spam_list", "../file/spam.txt", Docroot)
	InitnodeList = getRelativePathValue(i, "Path", "initnode_list", "../file/initnode.txt", Docroot)
	FollowList = getRelativePathValue(i, "Path", "follow_list", "../file/follow_list.txt", Docroot)
	NodeAllowFile = getRelativePathValue(i, "Path", "node_allow", "../file/node_allow.txt", Docroot)
	NodeDenyFile = getRelativePathValue(i, "Path", "node_deny", "../file/node_deny.txt", Docroot)
	ReAdminStr = getStringValue(i, "Gateway", "admin", "^(127|\\[::1\\])")
	ReFriendStr = getStringValue(i, "Gateway", "friend", "^(127|\\[::1\\])")
	ReVisitorStr = getStringValue(i, "Gateway", "visitor", ".")
	ServerName = getStringValue(i, "Gateway", "server_name", "")
	TagSize = getIntValue(i, "Gateway", "tag_size", 20)
	RSSRange = getInt64Value(i, "Gateway", "rss_range", 3*24*60*60)
	TopRecentRange = getInt64Value(i, "Gateway", "top_recent_range", 3*24*60*60)
	RecentRange = getInt64Value(i, "Gateway", "recent_range", 31*24*60*60)
	RecordLimit = getIntValue(i, "Gateway", "record_limit", 2048)
	Enable2ch = getBoolValue(i, "Gateway", "enable_2ch", false)
	EnableProf = getBoolValue(i, "Gateway", "enable_prof", false)
	HeavyMoon = getBoolValue(i, "Gateway", "moonlight", false)
	EnableEmbed = getBoolValue(i, "Gateway", "enable_embed", true)
	RelayNumber = getIntValue(i, "Gateway", "relay_number", 5)
	LogDir = getPathValue(i, "Path", "log_dir", "./log") //path from cwd
	ThreadPageSize = getIntValue(i, "Application Thread", "page_size", 50)
	DefaultThumbnailSize = getStringValue(i, "Application Thread", "thumbnail_size", "")
	ForceThumbnail = getBoolValue(i, "Application Thread", "force_thumbnail", false)
	ctype := "Application Thread"
	SaveRecord = getInt64Value(i, ctype, "save_record", 0)
	SaveSize = getIntValue(i, ctype, "save_size", 1)
	GetRange = getInt64Value(i, ctype, "get_range", 31*24*60*60)
	if GetRange > time.Now().Unix() {
		log.Fatal("get_range is too big")
	}
	SyncRange = getInt64Value(i, ctype, "sync_range", 10*24*60*60)
	if SyncRange > time.Now().Unix() {
		log.Fatal("sync_range is too big")
	}
	SaveRemoved = getInt64Value(i, ctype, "save_removed", 50*24*60*60)
	if SaveRemoved > time.Now().Unix() {
		log.Fatal("save_removed is too big")
	}

	if SyncRange == 0 {
		SaveRecord = 0
	}

	if SaveRemoved != 0 && SaveRemoved <= SyncRange {
		SyncRange = SyncRange + 1
	}

}

//Motd returns path to motd.txt
func Motd() string {
	return FileDir + "/motd.txt"
}

//Recent returns path to recent.txt
func Recent() string {
	return RunDir + "/recent.txt"
}

//AdminSid returns path to sid.txt
func AdminSid() string {
	return RunDir + "/sid.txt"
}

//PID returns path to pid.txt
func PID() string {
	return RunDir + "/pid.txt"
}

//Lookup returns path to lookup.txt
func Lookup() string {
	return RunDir + "/lookup.txt"
}

//Sugtag returns path to sugtag.txt
func Sugtag() string {
	return RunDir + "/sugtag.txt"
}

//Datakey returns path to datakey.txt
func Datakey() string {
	return RunDir + "/datakey.txt"
}