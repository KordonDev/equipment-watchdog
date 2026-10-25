#!/usr/bin/env bash
# Reserve glove IDs in bulk via POST /api/glove-ids/
#
# Usage:
#   ./scripts/add-glove-ids.sh <start> <end>     # range, inclusive
#   ./scripts/add-glove-ids.sh -l <id> [<id>...]  # explicit list
#
# Environment:
#   API_URL      (default: http://localhost:8080/api)
#   EW_USERNAME  (required)
#   EW_PASSWORD  (prompted if not set)

set -euo pipefail

API_URL="${API_URL:-http://localhost:8080/api}"

usage() {
    echo "Usage: $0 <start> <end>" >&2
    echo "       $0 -l <id> [<id>...]" >&2
    exit 1
}

json_escape() {
    local s="${1//\\/\\\\}"
    printf '%s' "${s//\"/\\\"}"
}

# Collect IDs from arguments
ids=()
if [ "${1:-}" = "-l" ]; then
    shift
    [ $# -ge 1 ] || usage
    ids=("$@")
else
    [ $# -eq 2 ] || usage
    [[ "$1" =~ ^[0-9]+$ && "$2" =~ ^[0-9]+$ ]] || usage
    [ "$1" -le "$2" ] || usage
    for ((i = $1; i <= $2; i++)); do
        ids+=("$i")
    done
fi

# Credentials
if [ -z "${EW_USERNAME:-}" ]; then
    echo "EW_USERNAME is not set" >&2
    exit 1
fi
if [ -z "${EW_PASSWORD:-}" ]; then
    read -r -s -p "Password for $EW_USERNAME: " EW_PASSWORD
    echo
fi

# Login and extract JWT from Set-Cookie header.
# The cookie is marked Secure, so we pass it manually instead of using a cookie jar.
login_body="{\"username\":\"$(json_escape "$EW_USERNAME")\",\"password\":\"$(json_escape "$EW_PASSWORD")\"}"
headers=$(curl -s -D - -o /dev/null -X POST "$API_URL/password-login" \
    -H 'Content-Type: application/json' \
    -d "$login_body") || { echo "Login request failed ($API_URL)" >&2; exit 1; }

TOKEN=$(printf '%s' "$headers" | tr -d '\r' | sed -n 's/^[Ss]et-[Cc]ookie: Authorization=\([^;]*\).*/\1/p' | head -n 1)
if [ -z "$TOKEN" ]; then
    echo "Login failed for user '$EW_USERNAME'" >&2
    exit 1
fi

# Create glove IDs
created=0
skipped=0
failed=0
for id in "${ids[@]}"; do
    code=$(curl -s -o /dev/null -w '%{http_code}' -X POST "$API_URL/glove-ids/" \
        -H 'Content-Type: application/json' \
        -H "Cookie: Authorization=$TOKEN" \
        -d "{\"gloveId\":\"$(json_escape "$id")\"}") || code="000"
    case "$code" in
        201) echo "✔ $id angelegt"; created=$((created + 1)) ;;
        409) echo "• $id existiert bereits"; skipped=$((skipped + 1)) ;;
        *)   echo "✘ $id Fehler (HTTP $code)"; failed=$((failed + 1)) ;;
    esac
done

echo
echo "Angelegt: $created, übersprungen: $skipped, fehlgeschlagen: $failed"
[ "$failed" -eq 0 ]
