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
