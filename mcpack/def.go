package mcpack

const (
	MCPACKV2_INVALID      = 0x0
	MCPACKV2_OBJECT       = 0x10
	MCPACKV2_ARRAY        = 0x20
	MCPACKV2_STRING       = 0x50
	MCPACKV2_RAW          = 0x60
	MCPACKV2_INT_8        = 0x11
	MCPACKV2_INT_16       = 0x12
	MCPACKV2_INT_32       = 0x14
	MCPACKV2_INT_64       = 0x18
	MCPACKV2_UINT_8       = 0x21
	MCPACKV2_UINT_16      = 0x22
	MCPACKV2_UINT_32      = 0x24
	MCPACKV2_UINT_64      = 0x28
	MCPACKV2_BOOL         = 0x31
	MCPACKV2_FLOAT        = 0x44
	MCPACKV2_DOUBLE       = 0x48
	MCPACKV2_DATE         = 0x58
	MCPACKV2_NULL         = 0x61
	MCPACKV2_SHORT_ITEM   = 0x80
	MCPACKV2_FIXED_ITEM   = 0xf
	MCPACKV2_DELETED_ITEM = 0x70
)
