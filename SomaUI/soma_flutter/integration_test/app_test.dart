import 'package:flutter/material.dart';
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

    // Test each example file in the IDE
    final exampleFiles = [
      'station_scaling.cljv',
      'grovers_search.cljv',
      'shors_factorization.cljv',
      'bell_state.cljv'
    ];

    for (final filename in exampleFiles) {
      testWidgets('verify live synthesis for $filename', (tester) async {
        app.main();
        await tester.pumpAndSettle();

        // Open IDE
        await tester.tap(find.text('OPEN CLOJUREV IDE'));
        await tester.pumpAndSettle();

        // Ensure IDE is open
        expect(find.text('ClojureV Forge'), findsOneWidget);

        // Tap the specific example file in the sidebar
        await tester.tap(find.text(filename));
        await tester.pumpAndSettle();

        // Find the synthesis button (play arrow icon) and tap it.
        final runButton = find.byIcon(Icons.play_arrow);
        expect(runButton, findsOneWidget);
        await tester.tap(runButton);
        
        // Wait for the synthesis API call to complete.
        // Since it relies on the local Go server which might take a second, we wait a bit longer.
        await tester.pump(const Duration(seconds: 3));
        await tester.pumpAndSettle();

        // Verify the success message appears in the terminal output
        expect(find.textContaining('[SUCCESS] Manifest complete'), findsWidgets);
      });
    }
  });
}
