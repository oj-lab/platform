package gravatar_utils

import (
	"crypto/md5"
	"fmt"
)

const gravatarURLFormat = "https://www.gravatar.com/avatar/%x?d=identicon&s=512"

func GetAvatarURL(identityStr string) string {
	md5sum := md5.Sum([]byte(identityStr))
	return fmt.Sprintf(gravatarURLFormat, md5sum)
}
