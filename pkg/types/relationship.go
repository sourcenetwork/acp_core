package types

func NewRelationship(resource, objId, relation, subjResource, subjId string) *Relationship {
	return &Relationship{
		Object: &Object{
			Resource: resource,
			Id:       objId,
		},
		Relation: relation,
		Subject: &Subject{
			Subject: &Subject_Object{
				Object: &Object{
					Resource: subjResource,
					Id:       subjId,
				},
			},
		},
	}
}

func NewActorRelationship(resource, objId, relation, actor string) *Relationship {
	return &Relationship{
		Object: &Object{
			Resource: resource,
			Id:       objId,
		},
		Relation: relation,
		Subject: &Subject{
			Subject: &Subject_Actor{
				Actor: &Actor{
					Id: actor,
				},
			},
		},
	}
}

func NewActorSetRelationship(resource, objId, relation, subjResource, subjId, subjRel string) *Relationship {
	return &Relationship{
		Object: &Object{
			Resource: resource,
			Id:       objId,
		},
		Relation: relation,
		Subject: &Subject{
			Subject: &Subject_ActorSet{
				ActorSet: &ActorSet{
					Object: &Object{
						Resource: subjResource,
						Id:       subjId,
					},
					Relation: subjRel,
				},
			},
		},
	}
}

func NewAllActorsRelationship(resource, objId, relation string) *Relationship {
	return &Relationship{
		Object: &Object{
			Resource: resource,
			Id:       objId,
		},
		Relation: relation,
		Subject: &Subject{
			Subject: &Subject_AllActors{
				AllActors: &AllActors{},
			},
		},
	}
}

func NewObject(resource, id string) *Object {
	return &Object{
		Resource: resource,
		Id:       id,
	}
}

func NewActor(actorId string) *Actor {
	return &Actor{
		Id: actorId,
	}
}
