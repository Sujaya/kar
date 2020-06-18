package test;

import com.ibm.research.kar.actor.ActorRef;
import com.ibm.research.kar.actor.exceptions.ActorMethodNotFoundException;

import static com.ibm.research.kar.Kar.actorCall;
import static com.ibm.research.kar.Kar.actorRef;

import javax.json.Json;
import javax.json.JsonObject;
import javax.json.JsonValue;

public class RunKar {

	public static void main(String[] args) {
	
		JsonObject params = Json.createObjectBuilder().add("number", 42).build();

		ActorRef dummy = actorRef("dummy", "dummyid");
		JsonValue result = null;
		try {
			result = actorCall(dummy, "canBeInvoked", params);
			
			System.out.println(result);
		} catch (ActorMethodNotFoundException e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}
	}

}
