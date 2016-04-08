
import java.io.File;
import java.io.FileInputStream;
import java.io.InputStream;
import java.io.PrintWriter;

import com.google.gson.Gson;
import com.google.gson.GsonBuilder;

import edu.cmu.sphinx.api.Configuration;
import edu.cmu.sphinx.api.SpeechResult;
import edu.cmu.sphinx.api.StreamSpeechRecognizer;
import edu.cmu.sphinx.result.WordResult;

public class TranscriberDemo {

	public static void main(String[] args) throws Exception {
		transcribe_file("files/wildshort");
		return;
	}

	public static void transcribe_file(String name) throws Exception {

		Configuration configuration = new Configuration();

		configuration.setAcousticModelPath("cmusphinx-en-us");
		configuration.setDictionaryPath("resource:/edu/cmu/sphinx/models/en-us/cmudict-en-us.dict");
		configuration.setLanguageModelPath("en-us.lm");

		StreamSpeechRecognizer recognizer = new StreamSpeechRecognizer(configuration);
		InputStream stream = new FileInputStream(new File(name + ".wav"));

		recognizer.startRecognition(stream);
		SpeechResult result;

		StringBuilder sb_text = new StringBuilder();
		StringBuilder sb_metadata = new StringBuilder();

		while ((result = recognizer.getResult()) != null) {
			String hyp = result.getHypothesis();
			System.out.format("Hypothesis: %s\n", hyp);
			System.out.println("List of recognized words and their times:");
			for (WordResult r : result.getWords()) {
				System.out.println(r);
				sb_metadata.append(r + ", ");
			}
			sb_text.append(hyp + ' ');
		}

		// Json
		Transcription target = new Transcription();
		target.textTranscription = sb_text.toString();
		target.metaData = sb_metadata.toString();

		Gson gson = new GsonBuilder().setPrettyPrinting().disableHtmlEscaping().create();
		String json = gson.toJson(target);

		// Printing json to file
		try (PrintWriter out = new PrintWriter(name + "-json.txt")) {
			out.println(json);
		}

		recognizer.stopRecognition();
	}
}
