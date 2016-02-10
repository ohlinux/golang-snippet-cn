package fetcher

import (
	. "config"
	. "db"
	"flag"
	"fmt"
	"github.com/go-xorm/xorm"
	"strconv"
	"strings"
	. "utils"
)

type Fetcher struct {
	conf     FetcherConfig
	db       *xorm.Engine
	api      Api
	filePath string
	baseName string
	tarUrl   string
	md5Url   string
	Md5Path  string
}

func NewFetcher(config *Config, db Api) *Fetcher {
	return &Fetcher{conf: config.Fetcher, db: db.DB, api: db}
}

func (f *Fetcher) exec(res *AppOriStatu) (done bool) {
	// md5check times
	times := 0
	for ; times < f.conf.CheckRetry; times++ {
		// download
		SimpleExec(fmt.Sprintf(f.conf.CurlFormat, f.tarUrl, res.LocalPath))
		res.Md5String = Sfile.Md5Sum(res.LocalPath)
		// check
		if f.md5Url != "" {
			SimpleExec(fmt.Sprintf(f.conf.CurlFormat, f.md5Url, f.Md5Path))
			if !Sfile.Md5Cmp(Sfile.Md5Read(f.Md5Path), res.Md5String) {
				continue
			}
		}
		break
	}
	done = times < f.conf.CheckRetry
	if done {
		tmp := &AppOriStatu{}
		exist, err := f.db.Where("origin_id=?").And("type=?").And("md5_string=?", res.OriginId, res.Type, res.Md5String).Get(tmp)
		if err == nil && exist {
			res = tmp
		}
	}
	return
}

func (f *Fetcher) preHttpDownload(p *PreFetch, res *AppOriStatu) (err error) {
	err = nil

	f.tarUrl = p.Name
	f.md5Url = ""
	f.Md5Path = ""

	if res.Version == "" {
		tmpv := &AppOriStatu{}
		exist, queryerr := f.db.Where("origin_id=?").And("type=?", res.OriginId, res.Type).Desc("create_at").Get(tmpv)
		if queryerr != nil {
			err = queryerr
			return
		}
		if !exist {
			return
		}
		tmpversion, strerr := strconv.Atoi(tmpv.Version)
		if strerr != nil {
			err = strerr
			return
		}
		res.Version = strconv.Itoa(tmpversion + 1)
	}

	res.LocalPath = fmt.Sprintf("%s/%s", f.filePath, f.baseName)

	return
}

func (f *Fetcher) preFtpDownload(p *PreFetch, res *AppOriStatu) (err error) {
	err = nil

	f.tarUrl = p.Name
	f.md5Url = ""
	f.Md5Path = ""

	if res.Version == "" {
		tmpv := &AppOriStatu{}
		exist, queryerr := f.db.Where("origin_id=?").And("type=?", res.OriginId, res.Type).Desc("create_at").Get(tmpv)
		if queryerr != nil {
			err = queryerr
			return
		}
		if !exist {
			return
		}
		tmpversion, strerr := strconv.Atoi(tmpv.Version)
		if strerr != nil {
			err = strerr
			return
		}
		res.Version = strconv.Itoa(tmpversion + 1)
	}

	res.LocalPath = fmt.Sprintf("%s/%s", f.filePath, f.baseName)

	return
}

func (f *Fetcher) preScmDownload(p *PreFetch, res *AppOriStatu) {
	urlFormat := f.conf.ScmTarFormat
	md5Format := f.conf.ScmMD5Format

	version := strings.Replace(res.Version, ".", "-", -1)
	f.tarUrl = fmt.Sprintf(urlFormat, p.Name, f.baseName, version, f.baseName, res.Version)
	f.md5Url = fmt.Sprintf(md5Format, p.Name, f.baseName, version, f.baseName, version)
	res.LocalPath = fmt.Sprintf("%s/%s_%s.tar.gz", f.filePath, f.baseName, res.Version)
	f.Md5Path = fmt.Sprintf("%s.md5", res.LocalPath)
	return
}

func (f *Fetcher) preDownload(p *PreFetch, res *AppOriStatu) {
	if res.Id == 0 {
		res.Id = int64(f.api.GetNextId(f.conf.DbTable))
	}
	res.Type = p.Type
	res.OriginId = p.OriginId
	res.Version = p.Version

	p.Name = strings.TrimSpace(p.Name)

	indextemp := strings.LastIndex(p.Name, "/") + 1
	if indextemp != -1 {
		f.baseName = Substr(p.Name, indextemp)
	} else {
		f.baseName = p.Name
	}

	f.filePath = fmt.Sprintf("%s/%02d/%04d/%06d/%08d", f.conf.DataDir, int(res.Id/1000000), int(res.Id/10000), int(res.Id/100), res.Id)

	return
}

func (f *Fetcher) checkLocal(p *PreFetch, res *AppOriStatu) (checked bool) {
	checked = false
	// check module and version in db
	exist, err := f.db.Where("origin_id=?").And("type=?").And("version=?", p.OriginId, p.Type, p.Version).Get(res)
	if err == nil && exist {
		// check file modified or not
		checked = Sfile.Md5CheckStr(res.LocalPath, res.Md5String)
		if !checked {
			f.db.Delete(&AppOriStatu{Type: res.Type, OriginId: res.OriginId})
		}
	}
	return
}

func (f *Fetcher) Fetch(p *PreFetch) (res *AppOriStatu) {
	res = &AppOriStatu{}
	checked := f.checkLocal(p, res)

	if checked {
		// change file time
		SimpleExec("touch", res.LocalPath)
	} else {
		// basic download prepare
		f.preDownload(p, res)
		// create download url
		if strings.ToLower(p.Source) == "scm" {
			f.preScmDownload(p, res)
		} else if strings.ToLower(p.Source) == "ftp" {
			f.preFtpDownload(p, res)
		} else if strings.ToLower(p.Source) == "http" {
			f.preHttpDownload(p, res)
		}
		// download tar
		checked = f.exec(res)
		if checked {
			f.db.Insert(res)
		}
	}

	return
}

func FetchPackage(p *PreFetch) (res *AppOriStatu) {
	configPath := "../../conf/config.toml"

	flag.StringVar(&configPath,
		"conf",
		"../../conf/config.toml", "Path of the TOML configuration of the eggs.conf ,default conf/eggs.conf .")
	flag.Parse()

	conf, err := ConfigFromFile(configPath)
	if err != nil {
		panic(err.Error())
	}
	DBApi := Api{}
	DBApi.InitDB(conf.DataBase)
	DBApi.InitSchema(new(PreFetch), new(AppOriStatu))
	fetcher := NewFetcher(conf, DBApi)

	res = fetcher.Fetch(p)
	return
}
