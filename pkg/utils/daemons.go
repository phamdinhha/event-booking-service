package utils

import "context"

type Deamon func()

type DeamonGenerator func(ctx context.Context) (Deamon, error)
