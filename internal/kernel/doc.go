// Package kernel holds the small value types shared by every bounded context:
// money, quantities, phone numbers, business dates, IDs and the business error
// type. It depends only on the standard library and google/uuid, so any layer
// (domain included) may import it.
package kernel
