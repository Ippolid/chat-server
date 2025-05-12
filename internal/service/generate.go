package service

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate minimock -i ChatServerService -o ./mocks/ -s "_minimock.go"
//go:generate minimock -i github.com/Ippolid/platform_libary/pkg/db.TxManager -o mocks/tx_manager_minimock.go -n TxManagerMock -p mocks
