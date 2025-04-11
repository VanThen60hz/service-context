package emailc

import _ "embed"

//go:embed templates/otp.html
var otpTemplate string
