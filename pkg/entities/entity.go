package entities

import "time"

type ToEntity interface {
	ToEntity(ownerId string) (Entity, error)
}

type Entity struct {
	Name         string    //Name/title of the entity
	EntityUrl    string    //Url to the entity (so it can be opened in a browser)
	ExternalId   string    //Id for the entity (specific to the 3rd party system)
	Type         string    //Type of the entity
	ContentUrl   string    //Url or path to the full content for the entity that was downloaded
	OwnerId      string    //Id of the user who owns the entity (user that connected the connector)
	LastModified time.Time //Last modified time in 3rd party system

	Data   any  //Data of entity
	IsFile bool //Is entity a file
}
