# build_pokedex_go


# Cache
Cache works but when doing map > map > mapb I would expect the mapb to be in cache but the problem is how the URL is done, and that is the key of the cache which in the first attempt is without the offset but when doing mapb has the offset, the data on both is the same but the keys are different so that edge case is not in the cache. I we look at the 3rd map command, it shows that data was cached (2nd map). 
Pokedex > map       
URL https://pokeapi.co/api/v2/location-area/ not in cache, making HTTP request...
canalave-city-area  
eterna-city-area    
...
Pokedex > map
URL https://pokeapi.co/api/v2/location-area/?offset=20&limit=20 not in cache, making HTTP request...
mt-coronet-1f-route-216
mt-coronet-1f-route-211
...
Pokedex > mapb
URL https://pokeapi.co/api/v2/location-area/?offset=0&limit=20 not in cache, making HTTP request...
canalave-city-area
eterna-city-area
...
Pokedex > map       
URL https://pokeapi.co/api/v2/location-area/?offset=20&limit=20 in cache
mt-coronet-1f-route-216
mt-coronet-1f-route-211




