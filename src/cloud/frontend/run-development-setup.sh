#!/bin/bash

VUE_APP_PROFILE="development-setup" npm run serve
echo -ne "\x1b[?25h"
echo -ne "\x1b[0m"