import 'dart:io';

import 'package:flutter/material.dart';
import 'package:path_provider/path_provider.dart';
import 'package:realm/realm.dart';

import 'models/product.dart';

void main() async {
  const String appId = "mba-pos-giwrk";

  WidgetsFlutterBinding.ensureInitialized();

  MyApp.allProductsRealm = await createRealm(appId, CollectionType.products);
  // MyApp.importantTasksRealm =
  //     await createRealm(appId, CollectionType.importantTasks);
  // MyApp.normalTasksRealm = await createRealm(appId, CollectionType.normalTasks);

  runApp(const MyApp());
}

enum CollectionType { products, importantTasks, normalTasks }

Future<Realm> createRealm(String appId, CollectionType collectionType) async {
  final appConfig = AppConfiguration(appId);
  final app = App(appConfig);
  final user = await app.logIn(Credentials.anonymous());

  final flxConfig = Configuration.flexibleSync(user, [Product.schema],
      path: await absolutePath("db_${collectionType.name}.realm"));
  var realm = Realm(flxConfig);
  print("Created local realm db at: ${realm.config.path}");

  final RealmResults<Product> query;
  if (collectionType == CollectionType.products) {
    query = realm.all<Product>();
  } else {
    query = realm.query<Product>(r'isImportant == $0',
        [collectionType == CollectionType.importantTasks]);
  }

  realm.subscriptions.update((mutableSubscriptions) {
    mutableSubscriptions.add(query);
  });

  await realm.subscriptions.waitForSynchronization();
  print("Syncronization completed for realm: ${realm.config.path}");
  return realm;
}

Future<String> absolutePath(String fileName) async {
  final appDocsDirectory = await getApplicationDocumentsDirectory();
  final realmDirectory = '${appDocsDirectory.path}/mongodb-realm';
  if (!Directory(realmDirectory).existsSync()) {
    await Directory(realmDirectory).create(recursive: true);
  }
  return "$realmDirectory/$fileName";
}

class MyApp extends StatelessWidget {
  static late Realm allProductsRealm;
  // static late Realm importantTasksRealm;
  // static late Realm normalTasksRealm;

  const MyApp({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Flutter Demo',
      theme: ThemeData(
        primarySwatch: Colors.green,
      ),
      home: const MyHomePage(title: 'Flutter Realm Flexible Sync'),
    );
  }
}

class MyHomePage extends StatefulWidget {
  const MyHomePage({Key? key, required this.title}) : super(key: key);

  final String title;

  @override
  State<MyHomePage> createState() => _MyHomePageState();
}

class _MyHomePageState extends State<MyHomePage> {
  int _allTasksCount = MyApp.allProductsRealm.all<Product>().length;
  // int _importantTasksCount = MyApp.importantTasksRealm.all<Product>().length;
  // int _normalTasksCount = MyApp.normalTasksRealm.all<Product>().length;

  void _createImportantTasks() async {
    // var importantTasks = MyApp.importantTasksRealm.all<Product>();
    print("Creating important tasks");
    var allTasksCount = MyApp.allProductsRealm.all<Product>();
    MyApp.allProductsRealm.write(() {
      MyApp.allProductsRealm.add(Product(ObjectId(), "p1", "prodtype1",
          "servType1", "searchName1", "name1", true));
    });
    await MyApp.allProductsRealm.syncSession.waitForUpload();
//    await MyApp.importantTasksRealm.subscriptions.waitForSynchronization();
    setState(() {
      //    _importantTasksCount = importantTasks.length;
      _allTasksCount = allTasksCount.length;
    });
  }

  void _createNormalTasks() async {
    // var normalTasks = MyApp.normalProductsRealm.all<Product>();
    var allTasksCount = MyApp.allProductsRealm.all<Product>();
    MyApp.allProductsRealm.write(() {
      MyApp.allProductsRealm.add(Product(ObjectId(), "p2", "prodtype2",
          "servType2", "searchName2", "name2", true));
    });
    await MyApp.allProductsRealm.syncSession.waitForUpload();
    //await MyApp.normalTasksRealm.subscriptions.waitForSynchronization();
    setState(() {
      //_normalTasksCount = normalTasks.length;
      _allTasksCount = allTasksCount.length;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        appBar: AppBar(
          title: Text(widget.title),
        ),
        body: Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: <Widget>[
              const Text('Important tasks count:',
                  style: TextStyle(fontWeight: FontWeight.bold)),
              // Text('$_importantTasksCount',
              //     style: Theme.of(context).textTheme.headline4),
              // Text('Realm path: ${MyApp.importantTasksRealm.config.path}'),
              // const Text('Normal tasks count:',
              //     style: TextStyle(fontWeight: FontWeight.bold)),
              // Text('$_normalTasksCount',
              //     style: Theme.of(context).textTheme.headline4),
              // Text('Realm path: ${MyApp.normalTasksRealm.config.path}'),
              // const Text('All tasks count:',
              //     style: TextStyle(fontWeight: FontWeight.bold)),
              Text('$_allTasksCount',
                  style: Theme.of(context).textTheme.headline4),
              Text('Realm path: ${MyApp.allProductsRealm.config.path}'),
            ],
          ),
        ),
        floatingActionButton: Stack(
          children: <Widget>[
            Padding(
              padding: const EdgeInsets.only(left: 0.0),
              child: Align(
                  alignment: Alignment.bottomLeft,
                  child: FloatingActionButton(
                    onPressed: _createImportantTasks,
                    tooltip: 'High priority task',
                    child: const Icon(Icons.add),
                  )),
            ),
            Padding(
              padding: const EdgeInsets.only(right: 40.0),
              child: Align(
                  alignment: Alignment.bottomRight,
                  child: FloatingActionButton(
                    onPressed: _createNormalTasks,
                    tooltip: 'Normal task',
                    child: const Icon(Icons.add),
                  )),
            ),
          ],
        ),
        floatingActionButtonLocation: FloatingActionButtonLocation.startFloat);
  }
}
