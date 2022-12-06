package historyManager

const (
	SELECT_HISTORY_ENTRIES_BY_ASSET_ID = `
		SELECT id, asset_id, price, direction, perc, date
		FROM history
		WHERE asset_id=$1
		ORDER BY date DESC
		LIMIT $2
		OFFSET $3;
	`
	SELECT_ASSETS = `
		SELECT id, name FROM assets;
	`
	INSERT_HISTORY_ENTRIES = `
		INSERT INTO history (id, asset_id, price, direction, perc, date)
		VALUES
	`
	INSERT_ASSET = `
		INSERT INTO assets (id, name)
		VALUES ($1, $2);
	`
)