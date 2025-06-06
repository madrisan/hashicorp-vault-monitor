package models
// Owner type of device.
type ManagedDeviceOwnerType int

const (
    // Unknown device owner type.
    UNKNOWN_MANAGEDDEVICEOWNERTYPE ManagedDeviceOwnerType = iota
    // Corporate device owner type.
    COMPANY_MANAGEDDEVICEOWNERTYPE
    // Personal device owner type.
    PERSONAL_MANAGEDDEVICEOWNERTYPE
    // Evolvable enumeration sentinel value. Do not use.
    UNKNOWNFUTUREVALUE_MANAGEDDEVICEOWNERTYPE
)

func (i ManagedDeviceOwnerType) String() string {
    return []string{"unknown", "company", "personal", "unknownFutureValue"}[i]
}
func ParseManagedDeviceOwnerType(v string) (any, error) {
    result := UNKNOWN_MANAGEDDEVICEOWNERTYPE
    switch v {
        case "unknown":
            result = UNKNOWN_MANAGEDDEVICEOWNERTYPE
        case "company":
            result = COMPANY_MANAGEDDEVICEOWNERTYPE
        case "personal":
            result = PERSONAL_MANAGEDDEVICEOWNERTYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_MANAGEDDEVICEOWNERTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeManagedDeviceOwnerType(values []ManagedDeviceOwnerType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ManagedDeviceOwnerType) isMultiValue() bool {
    return false
}
