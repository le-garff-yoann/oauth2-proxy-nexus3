import org.sonatype.nexus.capability.CapabilityReference
import org.sonatype.nexus.capability.CapabilityType
import org.sonatype.nexus.internal.capability.DefaultCapabilityReference
import org.sonatype.nexus.internal.capability.DefaultCapabilityRegistry

def capabilityRegistry = container.lookup(DefaultCapabilityRegistry.class.getName())
def capabilityType = CapabilityType.capabilityType("rutauth")

DefaultCapabilityReference capabilityExists = capabilityRegistry.all.find { CapabilityReference capabilityReference ->
    capabilityReference.context().descriptor().type() == capabilityType
}

if (capabilityExists) {
    return "Rut Auth capability is already configured"
} else {
    return sprintf("Rut Auth capability created as: %s", capabilityRegistry.add(
        capabilityType, true, "", [ "httpHeader": "X-Forwarded-User" ]
    ).toString())
}
