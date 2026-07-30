package application

import "errors"

var ErrNotPurchased = errors.New("solo puedes comentar en la publicación de un vendedor si le compraste algo")
