-- name: TransferAmount :one

INSERT INTO transfers (
    transfer_from_account, 
    transfer_to_account,   
    amount
) VALUES (
  $1, $2, $3
)RETURNING *;

-- name: GetTransfer :one
SELECT * FROM transfers
WHERE id = $1 
LIMIT 1;

-- name: ListTransfers :many
SELECT * FROM transfers
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: DeleteTransfers :exec
DELETE FROM transfers WHERE id = $1;
