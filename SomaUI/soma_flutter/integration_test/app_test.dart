import 'package:flutter_test/flutter_test.dart';
import 'package:integration_test/integration_test.dart';
import 'package:soma_flutter/main.dart' as app;

void main() {
  IntegrationTestWidgetsFlutterBinding.ensureInitialized();

  group('end-to-end test', () {
    testWidgets('verify dashboard layout and IDE toggle',
        (tester) async {
      app.main();
      await tester.pumpAndSettle();

      // Verify Dashboard Elements
      expect(find.text('SomaOS v3.5'), findsOneWidget);
      expect(find.text('XADC Thermal'), findsOneWidget);
      expect(find.text('OPEN CLOJUREV IDE'), findsOneWidget);

      // Open IDE
      await tester.tap(find.text('OPEN CLOJUREV IDE'));
      await tester.pumpAndSettle();

      // Verify IDE Elements
      expect(find.text('ClojureV Forge'), findsOneWidget);
    });
  });
}
