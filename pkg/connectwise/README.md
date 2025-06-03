## Example Usage
```
ctx := context.Background()

creds := connectwise.Creds{
    PublicKey:  "my_public_key",
    PrivateKey: "my_private_key",
    ClientId:   "my_client_id",
    CompanyId:  "my_company_id",
}

client := connectwise.NewClient(creds, http.DefaultClient)

params := &connectwise.QueryParams{
    OrderBy: "name asc",
    Conditions: "name='Help Desk'",
}

boards, err := client.ListBoards(ctx, params)
if err != nil {
    return nil, err
}
```