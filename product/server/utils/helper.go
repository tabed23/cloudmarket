package utils

import "go.mongodb.org/mongo-driver/bson/primitive"

// Helper function to remove duplicate ObjectIDs
func RemoveDuplicateObjectIDs(ids []primitive.ObjectID) []primitive.ObjectID {
    keys := make(map[primitive.ObjectID]bool)
    var result []primitive.ObjectID
    
    for _, id := range ids {
        if !keys[id] {
            keys[id] = true
            result = append(result, id)
        }
    }
    
    return result
}