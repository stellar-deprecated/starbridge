package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateDestinationAddressFnStellar_valid(t *testing.T) {
	assert.NoError(t, validateDestinationAddressFnStellar("GBJPEYVX2QHUXRWQH3SFWF63GDJQ426XD7QZRUOH5QKJ2ZXYDMKPMVX6"))
}

func TestValidateDestinationAddressFnStellar_invalid(t *testing.T) {
	assert.Error(t, validateDestinationAddressFnStellar("GBJPEYVX2QHUXRWQH3SFWF63GDJQ426XD7QZRUOH5QKJ2ZXYDMKPMVX"))
	assert.Error(t, validateDestinationAddressFnStellar("ABJPEYVX2QHUXRWQH3SFWF63GDJQ426XD7QZRUOH5QKJ2ZXYDMKPMVX6"))
	assert.Error(t, validateDestinationAddressFnStellar("hello"))
}
