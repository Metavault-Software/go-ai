#generations
curl https://api.openai.com/v1/images/generations \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $OPENAI_API_KEY" \
  -d '{
    "prompt": "a photo of a happy corgi puppy sitting and facing forward, studio light, longshot",
    "n":1,
    "size":"1024x1024"
   }'

#edits
#curl https://api.openai.com/v1/images/edits \
#  -H "Authorization: Bearer $OPENAI_API_KEY" \
#  -F image="@/Users/openai/happy_corgi.png" \
#  -F mask="@/Users/openai/mask.png" \
#  -F prompt="a photo of a happy corgi puppy with fancy sunglasses on sitting and facing forward, studio light, longshot" \
#  -F n=1 \
#  -F size="1024x1024"
#

##variations
#curl https://api.openai.com/v1/images/variations \
#  -H "Authorization: Bearer $OPENAI_API_KEY" \
#  -F image="@/Users/openai/corgi_with_sunglasses.png" \
#  -F n=4 \
#  -F size="1024x1024"
