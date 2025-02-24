package postgres

//func TestNewPool(t *testing.T) {
//	ctx := context.Background()
//	m, err := monitoring.New(monitoring.Config{})
//	require.NoError(t, err)
//
//	ctx = monitoring.SetInContext(ctx, m)
//	dbURL := os.Getenv("DATABASE_URL")
//
//	pool, err := NewPool(ctx, dbURL,
//		1, 1,
//		PoolMaxConnLifetime(1),
//		AttemptPingUponStartup(),
//	)
//	require.NoError(t, err)
//	require.NotNil(t, pool)
//}
