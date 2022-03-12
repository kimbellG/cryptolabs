package murmurhash3

func Sum32WithSeed(data []byte, seed uint32) uint32 {
	const (
		c1 = 0xcc9e2d51
		c2 = 0x1b873593
	)

	rotl32 := func(x, r uint32) uint32 {
		return (x << r) | (x >> (32 - r))
	}

	var (
		cur  = len(data)
		hash = seed
		k1   uint32
		i    = 0
	)

	for cur > 3 {
		k1, i = uint32(data[i])|
			uint32(data[i+1])<<8|
			uint32(data[i+2])<<16|
			uint32(data[i+3])<<24, i+4

		k1 *= c1
		k1 = rotl32(k1, 15)
		k1 *= c2

		hash ^= k1
		hash = rotl32(hash, 13)
		hash = hash*5 + 0xe6546b64

		cur -= 4
	}

	if cur&3 > 0 {
		switch cur {
		case 3:
			k1, i = uint32(data[i])|uint32(data[i+1])<<8|uint32(data[i+2])<<16, i+3
		case 2:
			k1, i = uint32(data[i])|uint32(data[i+1])<<8, i+2
		case 1:
			k1, i = uint32(data[i]), i+1
		}

		k1 *= c1
		k1 = rotl32(k1, 15)
		k1 *= c2
		hash ^= k1
	}

	hash ^= uint32(len(data))
	hash ^= hash >> 16
	hash *= 0x85ebca6b
	hash ^= hash >> 13
	hash *= 0xc2b2ae35
	hash ^= hash >> 16

	return hash
}
