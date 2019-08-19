package ipgeo

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/zhs007/ankadb"
	brucecorebase "github.com/zhs007/brucecore/base"
	ipgeopb "github.com/zhs007/brucecore/ipgeo/proto"
	"github.com/zhs007/jccclient"
	"go.uber.org/zap"
)

// DB -
type DB struct {
	ankaDB ankadb.AnkaDB
	client *jccclient.Client
}

// NewDB - new database
func NewDB(dbpath string, httpAddr string, engine string, client *jccclient.Client) (*DB, error) {
	cfg := ankadb.NewConfig()

	cfg.AddrHTTP = httpAddr
	cfg.PathDBRoot = dbpath
	cfg.ListDB = append(cfg.ListDB, ankadb.DBConfig{
		Name:   DBName,
		Engine: engine,
		PathDB: DBName,
	})

	ankaDB, err := ankadb.NewAnkaDB(cfg, nil)
	if ankaDB == nil {
		brucecorebase.Error("NewDB", zap.Error(err))

		return nil, err
	}

	brucecorebase.Info("NewDB", zap.String("dbpath", dbpath),
		zap.String("httpAddr", httpAddr), zap.String("engine", engine))

	db := &DB{
		ankaDB: ankaDB,
		client: client,
	}

	return db, err
}

// AddIPGeo - add a ipgeo
func (db *DB) AddIPGeo(ctx context.Context, ip string, geoip *ipgeopb.IPGeo) error {
	if ip == "" {
		return brucecorebase.ErrInvalidIP
	}

	if geoip == nil {
		return brucecorebase.ErrNoIPGeo
	}

	ipg, err := db.GetIPGeo(ctx, ip)
	if err == nil && ipg != nil {
		return nil
	}

	db.setIPGeo(ctx, ip, geoip)

	return nil
}

// GetIPGeo - get a ipgeo
func (db *DB) GetIPGeo(ctx context.Context, ip string) (*ipgeopb.IPGeo, error) {
	buf, err := db.ankaDB.Get(ctx, DBName, makeKey(ip))
	if err != nil {
		if err == ankadb.ErrNotFoundKey {
			return nil, nil
		}

		return nil, err
	}

	ri := &ipgeopb.IPGeo{}

	err = proto.Unmarshal(buf, ri)
	if err != nil {
		return nil, err
	}

	return ri, nil
}

// GetIPGeoEx - get a ipgeo
func (db *DB) GetIPGeoEx(ctx context.Context, ip string) (*ipgeopb.IPGeo, error) {
	if ip == "" {
		return nil, brucecorebase.ErrInvalidIP
	}

	buf, err := db.ankaDB.Get(ctx, DBName, makeKey(ip))
	if err != nil {
		if err == ankadb.ErrNotFoundKey {
			reply, err := db.client.GetGeoIP(ctx, ip, "")
			if err != nil {
				return nil, err
			}

			ipgeo := &ipgeopb.IPGeo{
				Latitude:     reply.Latitude,
				Longitude:    reply.Longitude,
				Organization: reply.Organization,
				Asn:          reply.Asn,
				Continent:    reply.Continent,
				Country:      reply.Country,
				Region:       reply.Region,
				City:         reply.City,
				Hostname:     reply.Hostname,
			}

			db.setIPGeo(ctx, ip, ipgeo)

			return ipgeo, nil
		}

		return nil, err
	}

	ri := &ipgeopb.IPGeo{}

	err = proto.Unmarshal(buf, ri)
	if err != nil {
		return nil, err
	}

	return ri, nil
}

// setIPGeo - set a ipgeo
func (db *DB) setIPGeo(ctx context.Context, ip string, geoip *ipgeopb.IPGeo) error {
	buf, err := proto.Marshal(geoip)
	if err != nil {
		return err
	}

	err = db.ankaDB.Set(ctx, DBName, makeKey(ip), buf)
	if err != nil {
		return err
	}

	return nil
}
