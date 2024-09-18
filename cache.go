package alertd

import (
	"context"
	"strconv"
	"time"

	"github.com/charmbracelet/log"
	"github.com/redis/go-redis/v9"
)

// Global variable for Redis client
var rdb *redis.Client = nil

// Connect to a Redis server using an `address`. This one has a format like
// `redis://<user>:<pass>@localhost:6379/<db>`.
// Set up the `rdb` global variable to that client. Returns an error if the
// passed address has a bad format.
func RedisConnect(address string) error {
	opt, err := redis.ParseURL(address)

	if err != nil {
		return err
	}

	rdb = redis.NewClient(opt)
	log.Infof("Connected to Redis server at `%s`", address)

	return nil
}

// Check if the server can use the already saved position for an user.
// This result depends on the movingActivity and timestamp values. Additionally,
// it checks if the cached lng/lat are the same of the input, in that case it
// can be skipped and server does not proceed to echo-request the GraphQL API.
// This cache system uses the user id as key of an hashset.
func UsePositionCache(input AlertInput) (*bool, error) {
	ctx := context.Background()
	var useCache bool = false

	// If there is no saved values for this key, just add the values and skip
	// the check for any timestamp
	hlen, err := rdb.HLen(ctx, input.Uid).Result()

	if err != nil {
		return nil, err
	}

	if hlen > 0 {
		pos, err := rdb.HGetAll(ctx, input.Uid).Result()

		if err != nil {
			return nil, err
		}

		timestamp, err := time.Parse("2006-01-02T15:04:05.999999999-07:00", pos["timestamp"])

		if err != nil {
			log.Warnf("Parsing timestamp: %s", err)
			timestamp, err = time.Parse("2006-01-02T15:04:05.999999999Z", pos["timestamp"])

			if err != nil {
				return nil, err
			}
		}

		diffTime := time.Now().Sub(timestamp).Seconds()

		lng, _ := strconv.ParseFloat(pos["longitude"], 64)
		lat, _ := strconv.ParseFloat(pos["latitude"], 64)
		useCache = lng == input.Longitude && lat == input.Latitude

		if !useCache {
			// An IN_VEHICLE values should be updated more quickly rather than STILL
			// values.
			switch pos["movingActivity"] {
			case "IN_VEHICLE":
				useCache = diffTime < 60
				break
			default:
				useCache = diffTime < 120
				break
			}
		}

		if useCache {
			log.Infof("Found a cached value for `%s` at time `%s`, with current diff time = `%f`s\nlng = %f, lng2 = %f, lat = %f, lat2 = %f", input.Uid, pos["timestamp"], diffTime, lng, input.Longitude, lat, input.Latitude)
		}
	}

	// Always updates the last position
	if err := rdb.HSet(
		ctx,
		input.Uid,
		map[string]interface{}{
			"longitude":      input.Longitude,
			"latitude":       input.Latitude,
			"movingActivity": input.MovingActivity,
			"timestamp":      time.Now(),
		}).Err(); err != nil {
		return nil, err
	}

	return &useCache, nil
}
