rmdir bin /s /q
rmdir temp /s /q
mkdir temp
mkdir bin
cd temp
git clone https://github.com/ipoluianov/gazer_client
git clone https://github.com/ipoluianov/gazer_node
cd gazer_node
go build -o ../../bin/gazer_node.exe ./main/main.go
cd ..
cd gazer_client
call flutter build windows
cd ..
cd ..
call ..\..\codesign\sing_gazer_node.bat
xcopy temp\gazer_client\build\windows\runner\Release bin\ /E
unzip temp\gazer_client\redist\redist.zip -d bin
call ..\..\codesign\sing_gazer_client.bat
install.nsi
call ..\..\codesign\sing_gazer_installer.bat
pause
