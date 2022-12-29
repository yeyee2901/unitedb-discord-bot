CREATE TABLE IF NOT EXISTS pokemon_battle_items (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
    `name` varchar(100) NOT NULL,
    `description` text NOT NULL,
    `tier` varchar(100) NOT NULL,
    `cooldown` varchar(100) NOT NULL,
    `trainer_level` int NOT NULL   
)
