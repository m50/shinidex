UPDATE pokemon 
SET shiny_locked = 0
WHERE id IN (
    'keldeo',
    'meloetta',
    'enamorus'
);

UPDATE pokemon_forms SET shiny_locked = 0
WHERE pokemon_id IN (
    'meloetta',
    'keldeo',
    'enamorus'
);