import QtQuick 2.2
import QtQuick.Window 2.1
import QtQuick.Controls 1.1
import QtQuick.Layouts 1.1

Window {
    title: "CS:GO Matchmaking Pinger"
    id: main_window
    visible: true
    minimumWidth: 300
    minimumHeight: 300
    width: 300
    height: 500

    /*
    ListModel {
       	id: serverModel
        objectName: server_list
        ListElement {
            description: "server"
        }
        ListElement {
            description: "status"
        }
        ListElement {
            description: "ping"
        }
    }
    */

    Item {
        
        TableView {
       	    model: servers.len
            id: table
            width: main_window.width
            height: main_window.height - ping_button.height - 10

		    TableViewColumn {
		 	    role: "name"
   			    title: "Server Name"
                delegate: Text {
                    text: servers.server(styleData.row).name
                }
			    width: 160
		    }
            TableViewColumn {
                role: "status"
                title: "Status"
                delegate: Text {
                    text: servers.server(styleData.row).status
                }
                width: 80
            }
		    TableViewColumn {
		 	    role: "ping"
			    title: "Ping"
                delegate: Text {
                    text: servers.server(styleData.row).ping
                }
			    width: 80
		    }

            Component.onCompleted: {
                servers.sync()
            }
        }

        Button {
            anchors { horizontalCenter: table.horizontalCenter;
                top: table.bottom; topMargin: 10 }
            id: ping_button
            onClicked: {
                ctrl.ping()
                servers.sync()
            }
       	    text: "Ping Servers"
        }
        
    }
}
