package application

import "errors"

var ErrNotPurchased = errors.New("solo puedes reseñar a un vendedor si le compraste algo")
