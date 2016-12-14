/**
 * Copyright (c) 2016 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package http

import (
	"net/http"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	NotFound      string = "not found"
	AlreadyExists string = "already exists"
	NotFoundEtcd  string = "cannnot get key"
	ConflictEtcd  string = "Compare failed"
)

func translateHttpErrorStatus(status int) int {
	if status == http.StatusNotFound {
		return http.StatusNotFound
	} else {
		return http.StatusInternalServerError
	}
}

func IsNotFoundError(err error) bool {
	return isErrorTypeStringInErrorMessage(NotFound, err) || isErrorTypeStringInErrorMessage(NotFoundEtcd, err)
}

func IsConflictError(err error) bool {
	return isErrorTypeStringInErrorMessage(ConflictEtcd, err)
}

func IsAlreadyExistsError(err error) bool {
	return isErrorTypeStringInErrorMessage(AlreadyExists, err)
}

func isErrorTypeStringInErrorMessage(errorType string, err error) bool {
	errorMessage := strings.ToUpper(err.Error())
	errorMessage = strings.TrimSpace(errorMessage)
	errorTypeString := strings.ToUpper(errorType)
	errorTypeStringLen := len(errorTypeString)

	index := 0
	for {
		index = strings.Index(errorMessage[index:], errorTypeString)
		if index == -1 {
			break
		}

		// assure errorTypeString string is not part of another word
		runeBefore, _ := utf8.DecodeLastRuneInString(errorMessage[:index])
		runeAfter, _ := utf8.DecodeRuneInString(errorMessage[index+errorTypeStringLen:])
		if (runeBefore == utf8.RuneError || unicode.IsSpace(runeBefore)) && (runeAfter == utf8.RuneError || unicode.IsSpace(runeAfter)) {
			return true
		}

		index += errorTypeStringLen
	}
	return false
}

func GetHttpStatusOrStatusError(status int, err error) int {
	if err != nil {
		if IsNotFoundError(err) {
			return http.StatusNotFound
		} else if IsAlreadyExistsError(err) {
			return http.StatusConflict
		} else if IsConflictError(err){
			return http.StatusConflict
		}
		return http.StatusInternalServerError
	}
	return status
}